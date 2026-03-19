import { useLoaderData } from "react-router-dom";

export const UserAnalytics = () => {
  const { totalMembers, activeUsersToday, recentEvents } = useLoaderData();

  return (
    <div className="p-8">
      <h1 className="text-2xl font-semibold text-white mb-8">Analytics</h1>

      {/* Top metrics */}
      <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-6 mb-10">

        {/* Team size */}
        <div className="bg-slate-800 border border-slate-700 rounded-xl p-6 shadow-lg hover:border-blue-500 hover:shadow-blue-500/10 transition duration-200">
          <div className="flex flex-col space-y-2">
            <span className="text-sm text-gray-400">
              Total Members
            </span>

            <span className="text-4xl font-bold text-white">
              {totalMembers}
            </span>
          </div>
        </div>

        {/* Active users today */}
        <div className="bg-slate-800 border border-slate-700 rounded-xl p-6 shadow-lg hover:border-blue-500 hover:shadow-blue-500/10 transition duration-200">
          <div className="flex flex-col space-y-2">
            <span className="text-sm text-gray-400">
              Active Users Today
            </span>

            <span className="text-4xl font-bold text-white">
              {activeUsersToday}
            </span>
          </div>
        </div>

      </div>

      {/* Activity Feed */}
      <div className="bg-slate-800 border border-slate-700 rounded-xl p-6 shadow-lg">
        <h2 className="text-lg font-semibold text-white mb-4">
          Recent Activity
        </h2>

        {recentEvents && recentEvents.length > 0 ? (
          <div className="space-y-4">
            {recentEvents.map((event, index) => (
              <div
                key={index}
                className="flex justify-between items-center border-b border-slate-700 pb-2"
              >
                <div className="flex flex-col">
                  <span className="text-sm text-white">
                    {event.event_type}
                  </span>
                  <span className="text-xs text-gray-400">
                    User: {event.user_id}
                  </span>
                </div>

                <span className="text-xs text-gray-500">
                  {new Date(event.event_time).toLocaleString()}
                </span>
              </div>
            ))}
          </div>
        ) : (
          <div className="text-gray-400 text-sm">
            No recent activity
          </div>
        )}
      </div>
    </div>
  );
};