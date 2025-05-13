import { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { useAuth } from "../AuthContext";

export default function FilesPage() {
  const { accessToken, logout } = useAuth();
  const [file, setFile] = useState(null);
  const [files, setFiles] = useState([]);
  const navigate = useNavigate();

  const fetchFiles = async () => {
    try {
      const res = await fetch("http://localhost:8080/api/files/list", {
        headers: {
          Authorization: `Bearer ${accessToken}`,
        },
      });

      if (!res.ok) {
        console.error("Ошибка при получении списка файлов", res.status);
        setFiles([]);
        return;
      }

      const data = await res.json();
      if (Array.isArray(data)) {
        setFiles(data.files);
      } else {
        console.warn("Сервер вернул не массив:", data);
        setFiles([]);
      }
    } catch (err) {
      console.error("Ошибка сети:", err);
      setFiles([]);
    }
  };

  const upload = async () => {
    if (!file) return;
    const formData = new FormData();
    formData.append("file", file);

    await fetch("http://localhost:8080/api/files/upload", {
      method: "POST",
      headers: {
        Authorization: `Bearer ${accessToken}`,
      },
      body: formData,
    });
    fetchFiles();
  };

  const download = async (filename) => {
    const res = await fetch(`http://localhost:8080/api/files/download/${filename}`, {
      headers: {
        Authorization: `Bearer ${accessToken}`,
      },
    });
    const blob = await res.blob();
    const a = document.createElement("a");
    a.href = window.URL.createObjectURL(blob);
    a.download = filename;
    a.click();
  };

  const remove = async (filename) => {
    await fetch(`http://localhost:8080/api/files/delete/${filename}`, {
      method: "DELETE",
      headers: {
        Authorization: `Bearer ${accessToken}`,
      },
    });
    fetchFiles();
  };

  const handleLogout = () => {
    logout();
    navigate("/login");
  };

  useEffect(() => {
    fetchFiles();
  }, []);

  return (
    <div className="max-w-2xl mx-auto py-10 px-4">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-semibold">Файлы</h1>
        <button onClick={handleLogout} className="text-sm text-red-500">Выйти</button>
      </div>

      <div className="flex gap-2 mb-4">
        <input type="file" onChange={(e) => setFile(e.target.files[0])} />
        <button onClick={upload} className="bg-blue-600 text-white px-4 py-1 rounded">Загрузить</button>
      </div>

      <ul className="space-y-2">
        {Array.isArray(files) && files.map((name) => (
          <li key={name} className="flex justify-between items-center border rounded px-4 py-2">
            <span>{name}</span>
            <div className="space-x-2">
              <button onClick={() => download(name)} className="text-blue-500">Скачать</button>
              <button onClick={() => remove(name)} className="text-red-500">Удалить</button>
            </div>
          </li>
        ))}
      </ul>
    </div>
  );
}
