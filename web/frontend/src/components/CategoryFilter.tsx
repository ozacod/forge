import type { Category } from '../types';

interface CategoryFilterProps {
  categories: Category[];
  selectedCategory: string | null;
  onSelect: (category: string | null) => void;
}

export function CategoryFilter({ categories, selectedCategory, onSelect }: CategoryFilterProps) {
  return (
    <div className="flex flex-wrap gap-2">
      <button
        onClick={() => onSelect(null)}
        className={`category-pill px-4 py-2 rounded-full text-sm font-medium transition-all ${
          selectedCategory === null ? 'active' : ''
        }`}
      >
        All
      </button>
      {categories.map((category) => (
        <button
          key={category.id}
          onClick={() => onSelect(category.id)}
          className={`category-pill px-4 py-2 rounded-full text-sm font-medium transition-all flex items-center gap-2 ${
            selectedCategory === category.id ? 'active' : ''
          }`}
        >
          <span>{category.icon}</span>
          <span>{category.name}</span>
        </button>
      ))}
    </div>
  );
}

