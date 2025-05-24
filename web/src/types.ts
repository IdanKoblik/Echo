export type Role = 'receiver' | 'sender';

export interface FormData {
  localPort: string;
  remoteAddr: string;
  filePath?: string;
  dest?: string;
}

export interface FileData {
  name: string;
  data: ArrayBuffer;
}