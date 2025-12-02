import type { ProjectConfig } from './types';

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
      cli_version: '0.0.44',
      name: 'cpx',
      description: 'C++ Project Generator',
    };
  }
  return response.json();
}

// Library/category endpoints removed - using vcpkg package names directly

export async function previewCMake(config: ProjectConfig): Promise<string> {
  const response = await fetch(`${API_BASE}/preview`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(config),
  });
  
  if (!response.ok) throw new Error('Preview failed');
  const data = await response.json();
  return data.cmake_content;
}

export async function generateProject(config: ProjectConfig): Promise<Blob> {
  const response = await fetch(`${API_BASE}/generate`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(config),
  });
  
  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.detail || 'Generation failed');
  }
  
  return response.blob();
}

// Recipe reload endpoint removed - using vcpkg now
