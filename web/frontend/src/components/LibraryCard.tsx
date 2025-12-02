import type { Library, LibrarySelection } from '../types';

interface LibraryCardProps {
  library: Library;
  selection: LibrarySelection | undefined;
  onToggle: (id: string) => void;
  onConfigureOptions: (library: Library) => void;
}

export function LibraryCard({ library, selection, onToggle, onConfigureOptions }: LibraryCardProps) {
  const isSelected = !!selection;
  const hasOptions = library.options && library.options.length > 0;
  const configuredOptionsCount = selection ? Object.keys(selection.options).length : 0;

  return (
    <div
      onClick={() => onToggle(library.id)}
      className={`
        library-card card-glass rounded-xl p-5 cursor-pointer
        transition-all duration-300 ease-out
        hover:scale-[1.02] hover:shadow-lg
        ${isSelected ? 'selected ring-1 ring-cyan-400/40' : ''}
      `}
    >
      <div className="flex items-start justify-between gap-3">
        <div className="flex-1 min-w-0">
          <div className="flex items-center gap-2 mb-2">
            <h3 className="font-display font-semibold text-white truncate">
              {library.name}
            </h3>
            {library.header_only && (
              <span className="px-2 py-0.5 text-[10px] font-mono uppercase tracking-wider bg-emerald-500/20 text-emerald-400 rounded-full border border-emerald-500/30">
                Header-only
              </span>
            )}
          </div>
          
          <p className="text-sm text-gray-400 line-clamp-2 mb-3">
            {library.description}
          </p>
          
          <div className="flex flex-wrap gap-1.5">
            {library.tags.slice(0, 3).map((tag) => (
              <span
                key={tag}
                className="tag px-2 py-0.5 text-[11px] font-mono text-gray-400 rounded-md"
              >
                {tag}
              </span>
            ))}
          </div>
        </div>
        
        <div className="flex flex-col items-end gap-2">
          <input
            type="checkbox"
            checked={isSelected}
            onChange={() => onToggle(library.id)}
            className="checkbox-custom"
            onClick={(e) => e.stopPropagation()}
          />
          <span className="text-[10px] font-mono text-gray-500">
            C++{library.cpp_standard}
          </span>
        </div>
      </div>
      
      <div className="mt-4 pt-3 border-t border-white/5 flex items-center justify-between">
        <a
          href={library.github_url}
          target="_blank"
          rel="noopener noreferrer"
          className="text-xs text-gray-500 hover:text-cyan-400 transition-colors flex items-center gap-1"
          onClick={(e) => e.stopPropagation()}
        >
          <svg className="w-3.5 h-3.5" fill="currentColor" viewBox="0 0 24 24">
            <path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/>
          </svg>
          GitHub
        </a>
        
        <div className="flex items-center gap-2">
          {hasOptions && (
            <button
              onClick={(e) => {
                e.stopPropagation();
                onConfigureOptions(library);
              }}
              className={`text-xs px-2 py-1 rounded-md transition-colors flex items-center gap-1 ${
                isSelected
                  ? 'bg-purple-500/20 text-purple-400 hover:bg-purple-500/30'
                  : 'bg-white/5 text-gray-500 hover:bg-white/10'
              }`}
            >
              <svg className="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
              </svg>
              Options
              {configuredOptionsCount > 0 && (
                <span className="ml-1 px-1.5 py-0.5 text-[10px] bg-purple-500/30 rounded-full">
                  {configuredOptionsCount}
                </span>
              )}
            </button>
          )}
          
          {library.alternatives.length > 0 && (
            <span className="text-[10px] text-gray-600">
              {library.alternatives.length} alt{library.alternatives.length > 1 ? 's' : ''}
            </span>
          )}
        </div>
      </div>
    </div>
  );
}
