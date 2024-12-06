const categories = [
  "All",
  "Conference",
  "Workshop",
  "Meetup",
  "Social",
  "Career Fair",
  "Graduation",
  "Other",
];

export const CategoryFilter = ({ selectedCategory, onCategorySelect }) => {
  return (
    <div className="max-w-7xl mx-auto my-8">
      <h2 className="text-xl font-semibold mb-4 text-purple-700">Categories</h2>
      <div className="flex gap-3 flex-wrap">
        {categories.map((category) => (
          <button
            key={category}
            onClick={() => onCategorySelect(category)}
            className={`px-4 py-2 rounded-full ${
              category === selectedCategory
                ? "bg-purple-600 text-white"
                : "bg-white text-gray-700 hover:bg-purple-50"
            } border border-purple-100 transition-colors`}
          >
            {category}
          </button>
        ))}
      </div>
    </div>
  );
};
