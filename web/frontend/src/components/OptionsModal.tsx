import { useState, useEffect } from 'react';
import type { Library, LibraryOption } from '../types';

interface OptionsModalProps {
  library: Library;
  currentOptions: Record<string, any>;
  onSave: (options: Record<string, any>) => void;
  onClose: () => void;
}

export function OptionsModal({ library, currentOptions, onSave, onClose }: OptionsModalProps) {
  const [options, setOptions] = useState<Record<string, any>>({});

  useEffect(() => {
    // Initialize with current options or defaults
    const initialOptions: Record<string, any> = {};
    for (const opt of library.options) {
      initialOptions[opt.id] = currentOptions[opt.id] ?? opt.default;
    }
    setOptions(initialOptions);
  }, [library, currentOptions]);

  const handleChange = (optionId: string, value: any) => {
    setOptions((prev) => ({ ...prev, [optionId]: value }));
  };

  const handleSave = () => {
    // Only save options that differ from defaults
    const changedOptions: Record<string, any> = {};
    for (const opt of library.options) {
      if (options[opt.id] !== opt.default) {
        changedOptions[opt.id] = options[opt.id];
      }
    }
    onSave(changedOptions);
  };

  const renderOption = (opt: LibraryOption) => {
    const value = options[opt.id];

    switch (opt.type) {
      case 'boolean':
        return (
          <label className="flex items-center justify-between cursor-pointer group">
            <div className="flex-1">
              <div className="text-sm font-medium text-white group-hover:text-cyan-400 transition-colors">
                {opt.name}
              </div>
              <div className="text-xs text-gray-500 mt-0.5">{opt.description}</div>
              {opt.cmake_var && (
                <div className="text-[10px] font-mono text-gray-600 mt-1">
                  CMake: {opt.cmake_var}
                </div>
              )}
            </div>
            <button
              type="button"
              onClick={() => handleChange(opt.id, !value)}
              className={`relative w-11 h-6 rounded-full transition-all ${
                value ? 'bg-cyan-500' : 'bg-white/10'
              }`}
            >
              <span
                className={`absolute top-1 w-4 h-4 bg-white rounded-full transition-all ${
                  value ? 'left-6' : 'left-1'
                }`}
              />
            </button>
          </label>
        );

      case 'string':
        return (
          <div>
            <label className="block text-sm font-medium text-white mb-1">
              {opt.name}
            </label>
            <div className="text-xs text-gray-500 mb-2">{opt.description}</div>
            <input
              type="text"
              value={value || ''}
              onChange={(e) => handleChange(opt.id, e.target.value)}
              placeholder={opt.default || 'Enter value...'}
              className="input-field w-full px-3 py-2 rounded-lg text-sm text-white font-mono"
            />
            {opt.cmake_var && (
              <div className="text-[10px] font-mono text-gray-600 mt-1">
                CMake: {opt.cmake_var}
              </div>
            )}
          </div>
        );

      case 'integer':
        return (
          <div>
            <label className="block text-sm font-medium text-white mb-1">
              {opt.name}
            </label>
            <div className="text-xs text-gray-500 mb-2">{opt.description}</div>
            <input
              type="number"
              value={value || 0}
              onChange={(e) => handleChange(opt.id, parseInt(e.target.value) || 0)}
              className="input-field w-full px-3 py-2 rounded-lg text-sm text-white font-mono"
            />
            {opt.cmake_var && (
              <div className="text-[10px] font-mono text-gray-600 mt-1">
                CMake: {opt.cmake_var}
              </div>
            )}
          </div>
        );

      case 'choice':
        return (
          <div>
            <label className="block text-sm font-medium text-white mb-1">
              {opt.name}
            </label>
            <div className="text-xs text-gray-500 mb-2">{opt.description}</div>
            <div className="flex flex-wrap gap-2">
              {opt.choices?.map((choice) => (
                <button
                  key={choice}
                  type="button"
                  onClick={() => handleChange(opt.id, choice)}
                  className={`px-3 py-1.5 rounded-lg text-sm font-mono transition-all ${
                    value === choice
                      ? 'bg-cyan-500/20 text-cyan-400 border border-cyan-500/40'
                      : 'bg-white/5 text-gray-400 border border-white/10 hover:bg-white/10'
                  }`}
                >
                  {choice}
                </button>
              ))}
            </div>
            {opt.cmake_var && (
              <div className="text-[10px] font-mono text-gray-600 mt-2">
                CMake: {opt.cmake_var}
              </div>
            )}
          </div>
        );

      default:
        return null;
    }
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
      {/* Backdrop */}
      <div
        className="absolute inset-0 bg-black/60 backdrop-blur-sm"
        onClick={onClose}
      />

      {/* Modal */}
      <div className="relative card-glass rounded-2xl w-full max-w-lg max-h-[80vh] flex flex-col animate-scale-in">
        {/* Header */}
        <div className="flex items-center justify-between p-5 border-b border-white/10">
          <div>
            <h2 className="font-display font-semibold text-lg text-white">
              {library.name} Options
            </h2>
            <p className="text-xs text-gray-500 mt-0.5">
              Configure build options for this library
            </p>
          </div>
          <button
            onClick={onClose}
            className="p-2 text-gray-500 hover:text-white transition-colors rounded-lg hover:bg-white/5"
          >
            <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        {/* Options list */}
        <div className="flex-1 overflow-y-auto p-5 space-y-5">
          {library.options.length === 0 ? (
            <div className="text-center py-8 text-gray-500">
              <p>No configurable options for this library.</p>
            </div>
          ) : (
            library.options.map((opt) => (
              <div
                key={opt.id}
                className="p-4 rounded-xl bg-white/5 border border-white/5"
              >
                {renderOption(opt)}
                {opt.affects_link && (
                  <div className="mt-2 text-[10px] text-amber-400/80 flex items-center gap-1">
                    <svg className="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                    </svg>
                    This option affects library linking
                  </div>
                )}
              </div>
            ))
          )}
        </div>

        {/* Footer */}
        <div className="flex items-center justify-end gap-3 p-5 border-t border-white/10">
          <button
            onClick={onClose}
            className="px-4 py-2 rounded-lg text-sm text-gray-400 hover:text-white hover:bg-white/5 transition-colors"
          >
            Cancel
          </button>
          <button
            onClick={handleSave}
            className="btn-primary px-4 py-2 rounded-lg text-sm font-medium text-white"
          >
            Save Options
          </button>
        </div>
      </div>
    </div>
  );
}

