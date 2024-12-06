import { Search, CalendarPlus } from "lucide-react";

export const Header = () => {
  return (
    <header className="bg-gradient-to-r from-purple-600 to-blue-600 p-4">
      <div className="max-w-7xl mx-auto flex items-center justify-between">
        <h1 className="text-white text-2xl font-bold">Event Hub</h1>

        <div className="flex-1 max-w-2xl mx-8">
          <div className="relative">
            <input
              type="text"
              placeholder="Search Events..."
              className="w-full px-4 py-2 rounded-lg bg-white/90 focus:outline-none focus:ring-2 focus:ring-purple-300"
            />
            <Search
              className="absolute right-3 top-2.5 text-gray-500"
              size={20}
            />
          </div>
        </div>

        <button className="bg-white text-purple-600 px-4 py-2 rounded-lg flex items-center gap-2 hover:bg-purple-50 transition-colors">
          <CalendarPlus size={20} />
          Create Event
        </button>
      </div>
    </header>
  );
};
