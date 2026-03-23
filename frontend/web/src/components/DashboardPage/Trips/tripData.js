export const tripRecords = [
  {
    id: "TRP-24031",
    routeName: "North Hub to City Center",
    driver: "Aarav Singh",
    vehicle: "MH12RX9044",
    origin: "North Distribution Hub",
    destination: "City Center Depot",
    status: "Created",
    date: "22 Mar 2026",
    notes: "Primary dispatch route for central retail deliveries.",
    location: {
      center: [25.2867, 51.5336],
      path: [
        [25.3012, 51.5166],
        [25.2964, 51.5231],
        [25.2909, 51.5284],
        [25.2867, 51.5336],
      ],
    },
  },
  {
    id: "TRP-24032",
    routeName: "Airport Corridor Loop",
    driver: "Fatima Noor",
    vehicle: "DL04KT1182",
    origin: "Airport Cargo Terminal",
    destination: "West Retail Cluster",
    status: "Created",
    date: "22 Mar 2026",
    notes: "Morning corridor trip with lighter commercial drop-offs.",
    location: {
      center: [25.2609, 51.5651],
      path: [
        [25.2735, 51.5628],
        [25.2679, 51.5664],
        [25.2631, 51.5696],
        [25.2609, 51.5651],
      ],
    },
  },
  {
    id: "TRP-24033",
    routeName: "Industrial Belt Express",
    driver: "Rohan Mehta",
    vehicle: "KA09VF2210",
    origin: "South Manufacturing Park",
    destination: "Riverfront Warehouse",
    status: "Started",
    date: "21 Mar 2026",
    notes: "Already started and moving on the industrial corridor.",
    location: {
      center: [25.2235, 51.5008],
      path: [
        [25.2124, 51.4899],
        [25.2179, 51.4947],
        [25.2215, 51.4972],
        [25.2235, 51.5008],
      ],
    },
  },
];

export const getTripById = (tripId) =>
  tripRecords.find((trip) => trip.id === tripId);
