import { useEffect, useState } from "react";
import { createMockTrackingStream } from "../utils/createMockTrackingStream";

function buildTrackingSocketUrl(channel, tripId) {
  const configuredBaseUrl = import.meta.env.VITE_BACKEND_WS_URL
    ?? import.meta.env.VITE_BACKEND_API_URL
    ?? window.location.origin;

  const socketUrl = new URL(configuredBaseUrl, window.location.origin);
  const currentPath = socketUrl.pathname.replace(/\/+$/, "");
  const apiBasePath = currentPath.endsWith("/v1") ? currentPath : `${currentPath}/v1`;

  socketUrl.protocol = socketUrl.protocol === "https:" ? "wss:" : "ws:";
  socketUrl.pathname = `${apiBasePath}/ws/${channel}-tracking`;

  if (tripId) {
    socketUrl.searchParams.set("tripId", tripId);
  }

  return socketUrl.toString();
}

export function useTrackingSocket({
  channel,
  tripId,
  driverId,
  vehicleId,
  routePoints,
}) {
  const [connectionState, setConnectionState] = useState("connecting");
  const [socketError, setSocketError] = useState(null);
  const [liveLocation, setLiveLocation] = useState(null);
  const [livePath, setLivePath] = useState([]);

  useEffect(() => {
    const socketUrl = buildTrackingSocketUrl(channel, tripId);
    const socket = new WebSocket(socketUrl);
    const getNextTrackingMessage = createMockTrackingStream({
      tripId,
      driverId,
      vehicleId,
      routePoints,
    });
    let intervalId;

    setLiveLocation(null);
    setLivePath([]);

    const sendNextLocation = () => {
      if (socket.readyState !== WebSocket.OPEN) {
        return;
      }

      const nextMessage = getNextTrackingMessage();
      const nextLocation = nextMessage.data;

      setLiveLocation(nextLocation);
      setLivePath((currentPath) => [...currentPath, nextLocation]);
      socket.send(JSON.stringify(nextMessage));
    };

    socket.onopen = () => {
      setConnectionState("connected");
      setSocketError(null);
      sendNextLocation();
      intervalId = window.setInterval(() => {
        sendNextLocation();
      }, 500);
    };

    socket.onerror = () => {
      setSocketError(`Unable to connect to the tracking websocket at ${socketUrl}.`);
    };

    socket.onclose = () => {
      setConnectionState("disconnected");
      if (intervalId) {
        window.clearInterval(intervalId);
      }
    };

    return () => {
      if (intervalId) {
        window.clearInterval(intervalId);
      }
      socket.close();
    };
  }, [channel, tripId, driverId, vehicleId, routePoints]);

  return {
    connectionState,
    liveLocation,
    livePath,
    socketError,
  };
}
