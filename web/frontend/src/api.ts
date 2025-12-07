// Use relative URL in production, localhost in development
const API_BASE = import.meta.env.DEV ? 'http://localhost:8000/api' : '/api';

export interface VersionInfo {
  version: string;
  cli_version: string;
  name: string;
  description: string;
}

export async function fetchVersion(): Promise<VersionInfo> {
  const response = await fetch(`${API_BASE}/version`);
  if (!response.ok) {
    // Return default version if server is not available
    return {
      version: '0.0.44',
      cli_version: '1.0.4',
      name: 'cpx',
      description: 'C++ Project Generator',
    };
  }
  return response.json();
}

