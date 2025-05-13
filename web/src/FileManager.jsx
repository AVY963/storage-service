import { useState, useEffect } from 'react';
import { Input } from './ui/Input';
import { Button } from './ui/Button';

export default function FileManager() {
  const [file, setFile] = useState(null);
  const [fileList, setFileList] = useState([]);
  const token = localStorage.getItem('token');

  const upload = async () => {
    if (!file) return;
    const formData = new FormData();
    formData.append('file', file);
    await fetch('http://localhost:8080/upload', {
      method: 'POST',
      headers: { Authorization: `Bearer ${token}` },
      body: formData,
    });
    listFiles();
  };

  const listFiles = async () => {
    const res = await fetch('http://localhost:8080/list', {
      headers: { Authorization: `Bearer ${token}` },
    });
    const data = await res.json();
    setFileList(data);
  };

  const download = async (filename) => {
    const res = await fetch(`http://localhost:8080/download/${filename}`, {
      headers: { Authorization: `Bearer ${token}` },
    });
    const blob = await res.blob();
    const a = document.createElement('a');
    a.href = window.URL.createObjectURL(blob);
    a.download = filename;
    a.click();
  };

  const deleteFile = async (filename) => {
    await fetch(`http://localhost:8080/delete/${filename}`, {
      method: 'DELETE',
      headers: { Authorization: `Bearer ${token}` },
    });
    listFiles();
  };

  useEffect(() => {
    if (token) listFiles();
  }, [token]);

  return (
    <div>
      <div className="space-y-4">
        <Input type="file" onChange={(e) => setFile(e.target.files[0])} />
        <Button onClick={upload}>Загрузить</Button>
        <ul className="list-disc list-inside mt-4">
          {fileList.map((f, i) => (
            <li key={i} className="flex justify-between">
              <span>{f}</span>
              <div className="space-x-2">
                <Button onClick={() => download(f)}>Скачать</Button>
                <Button variant="destructive" onClick={() => deleteFile(f)}>Удалить</Button>
              </div>
            </li>
          ))}
        </ul>
      </div>
    </div>
  );
}