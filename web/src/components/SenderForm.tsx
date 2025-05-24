import { useState, useRef } from 'react';
import { FormData, FileData } from '../types';
import { Upload, Edit3, FilePieChart } from 'lucide-react';

interface SenderFormProps {
  onSubmit: (data: FormData, success: boolean) => void;
}

export default function SenderForm({ onSubmit }: SenderFormProps) {
  const [formData, setFormData] = useState<FormData>({
    localPort: '',
    remoteAddr: '',
    filePath: '',
  });
  const wsRef = useRef<WebSocket | null>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    const ws = new WebSocket(`ws://0.0.0.0:8080/ws`);
    wsRef.current = ws;

    ws.onopen = async () => {
      const request = {
        mode: "send",
        port: formData.localPort,
        remote: formData.remoteAddr,
        file: formData.filePath, 
      };
      ws.send(JSON.stringify(request));
      console.log("Request sent:", request);
    };

    ws.onerror = (error) => {
      console.error('WebSocket error:', error);
      onSubmit(formData, false);
      alert('Failed to connect to the receiver');
    };

    onSubmit(formData, true);
  };

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setFormData(prev => ({ ...prev, [name]: value }));
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-4 w-full max-w-md">
      <div className="space-y-2">
        <label htmlFor="localPort" className="block text-sm font-medium dark:text-gray-200">
          Local Port
        </label>
        <input
          type="text"
          id="localPort"
          name="localPort"
          value={formData.localPort}
          onChange={handleChange}
          placeholder="e.g., 8080"
          className="w-full px-4 py-2 rounded-lg border dark:border-gray-600 bg-white dark:bg-gray-800
                    focus:outline-none focus:ring-2 focus:ring-indigo-500 dark:text-white
                    transition-colors duration-200"
          required
        />
      </div>

      <div className="space-y-2">
        <label htmlFor="remoteAddr" className="block text-sm font-medium dark:text-gray-200">
          Remote Address
        </label>
        <input
          type="text"
          id="remoteAddr"
          name="remoteAddr"
          value={formData.remoteAddr}
          onChange={handleChange}
          placeholder="e.g., 192.168.1.100:8080"
          className="w-full px-4 py-2 rounded-lg border dark:border-gray-600 bg-white dark:bg-gray-800
                    focus:outline-none focus:ring-2 focus:ring-indigo-500 dark:text-white
                    transition-colors duration-200"
          required
        />
      </div>

      <div className="space-y-2">
        <label className="block text-sm font-medium dark:text-gray-200">
          File Selection Method
        </label>

        <input
            type="text"
            name="filePath"
            value={formData.filePath}
            onChange={handleChange}
            placeholder="/path/to/your/file"
            className="w-full px-4 py-2 rounded-lg border dark:border-gray-600 bg-white dark:bg-gray-800
                      focus:outline-none focus:ring-2 focus:ring-indigo-500 dark:text-white
                      transition-colors duration-200"
            required
          />
      </div>

      <button
        type="submit"
        className="w-full py-3 px-6 bg-indigo-600 text-white rounded-lg font-medium
                  hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500
                  transition-colors duration-200 transform hover:scale-[1.02]"
      >
        Send File
      </button>
    </form>
  );
}