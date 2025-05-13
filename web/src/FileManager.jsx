import { useState, useEffect } from 'react';
import { useAuth } from './AuthContext';
import { 
  Box, 
  Typography, 
  Paper, 
  Button, 
  List, 
  ListItem, 
  ListItemText, 
  ListItemSecondaryAction, 
  IconButton, 
  Divider, 
  Alert, 
  CircularProgress,
  Snackbar,
  Grid,
  Container,
  AppBar,
  Toolbar,
  LinearProgress,
  Card,
  CardContent,
  CardHeader,
  Avatar,
  Tooltip,
  Fade,
  Zoom,
  Slide,
  Dialog,
  DialogActions,
  DialogContent,
  DialogContentText,
  DialogTitle
} from '@mui/material';
import { 
  CloudUpload, 
  Delete, 
  Download, 
  Refresh, 
  InsertDriveFile,
  Image,
  PictureAsPdf,
  Code,
  VideoFile,
  AudioFile,
  Archive,
  Description
} from '@mui/icons-material';

export default function FileManager() {
  const { accessToken } = useAuth();
  const [file, setFile] = useState(null);
  const [fileList, setFileList] = useState([]);
  const [loading, setLoading] = useState(false);
  const [loadingFiles, setLoadingFiles] = useState(false);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');
  const [filePreview, setFilePreview] = useState('');
  const [uploadProgress, setUploadProgress] = useState(0);
  const [uploading, setUploading] = useState(false);
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
  const [fileToDelete, setFileToDelete] = useState(null);

  console.log("FileManager rendered, accessToken:", accessToken ? "присутствует" : "отсутствует");

  const handleFileChange = (e) => {
    const selectedFile = e.target.files[0];
    if (selectedFile) {
      setFile(selectedFile);
      // Создаем превью для файла, если это изображение
      if (selectedFile.type.startsWith('image/')) {
        const reader = new FileReader();
        reader.onload = (e) => setFilePreview(e.target.result);
        reader.readAsDataURL(selectedFile);
      } else {
        setFilePreview('');
      }
    }
  };

  const upload = async () => {
    if (!file) {
      setError('Пожалуйста, выберите файл для загрузки');
      return;
    }
    
    try {
      setUploading(true);
      setUploadProgress(0);
      
      const formData = new FormData();
      formData.append('file', file);
      
      // Используем XMLHttpRequest для отслеживания прогресса загрузки
      const xhr = new XMLHttpRequest();
      xhr.open('POST', 'http://localhost:8080/api/files/upload');
      xhr.setRequestHeader('Authorization', `Bearer ${accessToken}`);
      
      xhr.upload.addEventListener('progress', (event) => {
        if (event.lengthComputable) {
          const progress = Math.round((event.loaded / event.total) * 100);
          setUploadProgress(progress);
        }
      });
      
      xhr.onload = () => {
        console.log("Статус ответа загрузки:", xhr.status);
        console.log("Ответ загрузки:", xhr.responseText);
        
        if (xhr.status >= 200 && xhr.status < 300) {
          try {
            const response = JSON.parse(xhr.responseText);
            
            if (response.success) {
              const successMessage = response.message || 'Файл успешно загружен';
              setSuccess(successMessage);
              console.log("Файл успешно загружен:", file.name);
            } else {
              setError('Ошибка при загрузке файла: ' + (response.error || 'Неизвестная ошибка'));
              console.error("Ошибка загрузки:", response.error);
            }
          } catch (e) {
            console.error("Ошибка при разборе ответа:", e);
            setSuccess('Файл загружен, но возникла ошибка при обработке ответа');
          }
          
          setFile(null);
          setFilePreview('');
          listFiles();
        } else {
          try {
            const errorResponse = JSON.parse(xhr.responseText);
            setError(`Ошибка при загрузке файла: ${errorResponse.error || xhr.status}`);
            console.error("Ошибка загрузки:", errorResponse.error);
          } catch (e) {
            setError(`Ошибка при загрузке файла: ${xhr.status}`);
            console.error("Ошибка при загрузке файла:", xhr.status);
          }
        }
        setUploading(false);
      };
      
      xhr.onerror = (e) => {
        console.error("Ошибка соединения с сервером:", e);
        setError('Ошибка соединения с сервером');
        setUploading(false);
      };
      
      xhr.send(formData);
    } catch (err) {
      console.error("Неожиданная ошибка при загрузке:", err);
      setError(err.message || 'Произошла неизвестная ошибка');
      setUploading(false);
    }
  };

  const listFiles = async () => {
    try {
      setLoadingFiles(true);
      console.log("Отправка запроса на получение списка файлов...");
      const res = await fetch('http://localhost:8080/api/files/list', {
        headers: { Authorization: `Bearer ${accessToken}` },
      });
      
      console.log("Статус ответа API:", res.status);
      
      if (!res.ok) {
        throw new Error(`Не удалось получить список файлов: ${res.status}`);
      }
      
      try {
        const data = await res.json();
        console.log("Ответ сервера (JSON):", data);
        
        let filesList = [];
        
        // Формат ответа API: { "files": [ { "name": "filename", "created_at": "date", "updated_at": "date" }, ... ] }
        if (data && data.files && Array.isArray(data.files)) {
          console.log("Получены данные в формате {files: [...]}");
          // Извлекаем только имена файлов
          filesList = data.files.map(file => file.name).filter(name => name && typeof name === 'string');
        } 
        // На случай, если API вернет прямой список файлов
        else if (data && Array.isArray(data)) {
          console.log("Получены данные в формате массива");
          filesList = data
            .filter(item => item && typeof item === 'string' || (typeof item === 'object' && item.name))
            .map(item => typeof item === 'string' ? item : item.name);
        } 
        // На случай, если API вернет список в другом поле
        else if (data && typeof data === 'object') {
          console.log("Получены данные в формате объекта, ищем массив внутри");
          // Проверяем все поля объекта на наличие массива
          for (const key in data) {
            if (Array.isArray(data[key])) {
              console.log(`Найден массив в поле ${key}`);
              filesList = data[key]
                .filter(item => item && (typeof item === 'string' || (typeof item === 'object' && item.name)))
                .map(item => typeof item === 'string' ? item : item.name);
              break;
            }
          }
        }
        
        console.log("Итоговый список файлов:", filesList);
        setFileList(filesList);
        
        if (filesList.length === 0) {
          console.log("Список файлов пуст");
        }
      } catch (err) {
        console.error("Ошибка при обработке ответа:", err);
        setError("Ошибка при обработке ответа сервера: " + err.message);
        setFileList([]);
      }
    } catch (err) {
      console.error("Ошибка сети:", err);
      setError(err.message);
      setFileList([]);
    } finally {
      setLoadingFiles(false);
    }
  };

  const download = async (filename) => {
    if (!filename) {
      setError('Имя файла не указано');
      return;
    }
    
    try {
      setLoading(true);
      console.log(`Скачивание файла: ${filename}`);
      
      const res = await fetch(`http://localhost:8080/api/files/download/${encodeURIComponent(filename)}`, {
        headers: { Authorization: `Bearer ${accessToken}` },
      });
      
      console.log("Статус ответа скачивания:", res.status);
      
      if (!res.ok) {
        let errorMessage = `Ошибка при скачивании файла: ${res.status}`;
        try {
          const errorData = await res.json();
          errorMessage = errorData.error || errorMessage;
        } catch (e) {
          // Если ответ не в формате JSON, используем стандартное сообщение
        }
        
        throw new Error(errorMessage);
      }
      
      const blob = await res.blob();
      const url = window.URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = filename;
      document.body.appendChild(a);
      a.click();
      document.body.removeChild(a);
      window.URL.revokeObjectURL(url);
      
      console.log(`Файл ${filename} успешно скачан`);
      setSuccess('Файл успешно скачан');
    } catch (err) {
      console.error("Ошибка при скачивании:", err);
      setError(err.message || 'Не удалось скачать файл');
    } finally {
      setLoading(false);
    }
  };

  const confirmDelete = (filename) => {
    if (!filename) {
      setError('Имя файла не указано');
      return;
    }
    setFileToDelete(filename);
    setDeleteDialogOpen(true);
  };

  const handleDeleteConfirm = async () => {
    if (fileToDelete) {
      await deleteFile(fileToDelete);
    }
    setDeleteDialogOpen(false);
    setFileToDelete(null);
  };

  const handleDeleteCancel = () => {
    setDeleteDialogOpen(false);
    setFileToDelete(null);
  };

  const deleteFile = async (filename) => {
    if (!filename) {
      setError('Имя файла не указано');
      return;
    }
    
    try {
      setLoading(true);
      console.log(`Удаление файла: ${filename}`);
      
      const res = await fetch(`http://localhost:8080/api/files/delete/${encodeURIComponent(filename)}`, {
        method: 'DELETE',
        headers: { Authorization: `Bearer ${accessToken}` },
      });
      
      console.log("Статус ответа удаления:", res.status);
      
      if (!res.ok) {
        let errorMessage = `Ошибка при удалении файла: ${res.status}`;
        try {
          const errorData = await res.json();
          errorMessage = errorData.error || errorMessage;
        } catch (e) {
          // Если ответ не в формате JSON, используем стандартное сообщение
        }
        
        throw new Error(errorMessage);
      }
      
      try {
        const response = await res.json();
        const successMessage = response.message || 'Файл успешно удален';
        console.log(`Файл ${filename} успешно удален:`, response);
        setSuccess(successMessage);
      } catch (e) {
        // Если ответ не в формате JSON, используем стандартное сообщение об успехе
        console.log(`Файл ${filename} успешно удален`);
        setSuccess('Файл успешно удален');
      }
      
      // Обновляем список файлов
      listFiles();
    } catch (err) {
      console.error("Ошибка при удалении:", err);
      setError(err.message || 'Не удалось удалить файл');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    if (accessToken) listFiles();
  }, [accessToken]);

  const handleCloseSnackbar = () => {
    setError('');
    setSuccess('');
  };

  // Получение иконки в зависимости от типа файла
  const getFileIcon = (filename) => {
    if (!filename) return <Description />;
    
    const extension = filename.split('.').pop().toLowerCase();
    
    switch (extension) {
      case 'jpg':
      case 'jpeg':
      case 'png':
      case 'gif':
      case 'webp':
        return <Image color="primary" />;
      case 'pdf':
        return <PictureAsPdf color="error" />;
      case 'zip':
      case 'rar':
      case '7z':
      case 'tar':
      case 'gz':
        return <Archive color="warning" />;
      case 'mp4':
      case 'avi':
      case 'mov':
      case 'mkv':
        return <VideoFile color="secondary" />;
      case 'mp3':
      case 'wav':
      case 'ogg':
        return <AudioFile color="success" />;
      case 'js':
      case 'jsx':
      case 'ts':
      case 'tsx':
      case 'html':
      case 'css':
      case 'py':
      case 'java':
      case 'c':
      case 'cpp':
        return <Code color="info" />;
      default:
        return <InsertDriveFile color="primary" />;
    }
  };

  // Форматирование размера файла
  const formatFileSize = (bytes) => {
    if (bytes === 0) return '0 Байт';
    
    const units = ['Байт', 'КБ', 'МБ', 'ГБ', 'ТБ'];
    const i = Math.floor(Math.log(bytes) / Math.log(1024));
    
    return `${(bytes / Math.pow(1024, i)).toFixed(2)} ${units[i]}`;
  };

  return (
    <Box sx={{ minHeight: '100vh', backgroundColor: '#f5f7fa' }}>
      <AppBar position="static" elevation={0} sx={{ backgroundColor: '#1976d2' }}>
        <Toolbar>
          <Typography variant="h5" component="h1" sx={{ fontWeight: 'bold' }}>
            Файловый менеджер
          </Typography>
          <Box sx={{ flexGrow: 1 }} />
          <Tooltip title="Обновить список файлов">
            <IconButton color="inherit" onClick={listFiles} disabled={loadingFiles}>
              <Refresh />
            </IconButton>
          </Tooltip>
        </Toolbar>
      </AppBar>
      
      <Container maxWidth="lg" sx={{ py: 4 }}>
        <Grid container spacing={3}>
          {/* Секция загрузки файла */}
          <Grid item xs={12} md={5}>
            <Zoom in={true} timeout={500}>
              <Paper 
                elevation={3} 
                sx={{ 
                  p: 3, 
                  height: '100%', 
                  borderRadius: 2,
                  transition: 'all 0.3s',
                  '&:hover': {
                    boxShadow: 6
                  }
                }}
              >
                <Typography 
                  variant="h6" 
                  gutterBottom 
                  sx={{ 
                    color: 'primary.main', 
                    fontWeight: 'bold',
                    display: 'flex',
                    alignItems: 'center',
                    mb: 2
                  }}
                >
                  <CloudUpload sx={{ mr: 1 }} /> Загрузка файла
                </Typography>
                <Divider sx={{ mb: 3 }} />
                
                <Box sx={{ mb: 3, display: 'flex', flexDirection: 'column', alignItems: 'center' }}>
                  <input
                    type="file"
                    id="fileInput"
                    style={{ display: 'none' }}
                    onChange={handleFileChange}
                  />
                  <label htmlFor="fileInput">
                    <Button
                      variant="outlined"
                      component="span"
                      startIcon={<CloudUpload />}
                      sx={{ 
                        mb: 2,
                        px: 3,
                        py: 1.2,
                        borderRadius: 2,
                        borderWidth: 2,
                        transition: 'all 0.2s',
                        '&:hover': {
                          borderWidth: 2,
                          transform: 'translateY(-2px)'
                        }
                      }}
                    >
                      Выберите файл
                    </Button>
                  </label>
                  
                  {file && (
                    <Fade in={!!file}>
                      <Card sx={{ width: '100%', mt: 2, borderRadius: 2 }}>
                        <CardHeader
                          avatar={
                            <Avatar sx={{ bgcolor: 'primary.main' }}>
                              {getFileIcon(file.name)}
                            </Avatar>
                          }
                          title={file.name}
                          subheader={formatFileSize(file.size)}
                        />
                        
                        {filePreview && (
                          <Box 
                            component="img" 
                            src={filePreview} 
                            sx={{ 
                              width: '100%',
                              maxHeight: '200px',
                              objectFit: 'contain',
                              borderRadius: '0 0 8px 8px'
                            }}
                            alt={file.name}
                          />
                        )}
                      </Card>
                    </Fade>
                  )}
                </Box>
                
                {uploading && (
                  <Box sx={{ width: '100%', mb: 2 }}>
                    <LinearProgress variant="determinate" value={uploadProgress} />
                    <Typography variant="body2" align="center" sx={{ mt: 1 }}>
                      {`${uploadProgress}%`}
                    </Typography>
                  </Box>
                )}
                
                <Box sx={{ display: 'flex', justifyContent: 'center' }}>
                  <Button
                    variant="contained"
                    color="primary"
                    onClick={upload}
                    disabled={!file || uploading}
                    startIcon={uploading ? <CircularProgress size={20} color="inherit" /> : <CloudUpload />}
                    sx={{ 
                      px: 4, 
                      py: 1.2,
                      borderRadius: 2,
                      transition: 'all 0.2s',
                      '&:not(:disabled):hover': {
                        transform: 'translateY(-2px)',
                        boxShadow: 6
                      }
                    }}
                  >
                    {uploading ? 'Загрузка...' : 'Загрузить'}
                  </Button>
                </Box>
              </Paper>
            </Zoom>
          </Grid>
          
          {/* Секция списка файлов */}
          <Grid item xs={12} md={7}>
            <Zoom in={true} timeout={700}>
              <Paper 
                elevation={3} 
                sx={{ 
                  p: 3, 
                  height: '100%',
                  borderRadius: 2,
                  transition: 'all 0.3s',
                  '&:hover': {
                    boxShadow: 6
                  }
                }}
              >
                <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
                  <Typography 
                    variant="h6" 
                    sx={{ 
                      color: 'primary.main', 
                      fontWeight: 'bold',
                      display: 'flex',
                      alignItems: 'center'
                    }}
                  >
                    <InsertDriveFile sx={{ mr: 1 }} /> Список файлов
                  </Typography>
                  <Button 
                    startIcon={<Refresh />} 
                    variant="text" 
                    color="primary" 
                    onClick={listFiles} 
                    disabled={loadingFiles}
                    size="small"
                  >
                    Обновить
                  </Button>
                </Box>
                <Divider sx={{ mb: 3 }} />
                
                {loadingFiles && (
                  <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', my: 5 }}>
                    <CircularProgress />
                  </Box>
                )}
                
                {!loadingFiles && fileList.length === 0 && (
                  <Box sx={{ py: 5, display: 'flex', flexDirection: 'column', alignItems: 'center' }}>
                    <InsertDriveFile sx={{ fontSize: 60, color: 'text.disabled', mb: 2 }} />
                    <Typography variant="h6" color="text.secondary">
                      Файлы отсутствуют
                    </Typography>
                    <Typography variant="body2" color="text.disabled" align="center" sx={{ mt: 1 }}>
                      Загрузите файлы, используя форму слева
                    </Typography>
                  </Box>
                )}
                
                {!loadingFiles && fileList.length > 0 && (
                  <List sx={{ 
                    maxHeight: 400, 
                    overflow: 'auto',
                    '&::-webkit-scrollbar': {
                      width: '8px',
                    },
                    '&::-webkit-scrollbar-track': {
                      backgroundColor: '#f1f1f1',
                      borderRadius: '10px',
                    },
                    '&::-webkit-scrollbar-thumb': {
                      backgroundColor: '#c1c1c1',
                      borderRadius: '10px',
                      '&:hover': {
                        backgroundColor: '#a8a8a8',
                      },
                    },
                  }}>
                    {fileList.map((filename, index) => (
                      <Fade key={index} in={true} timeout={200 + index * 100}>
                        <Box>
                          <ListItem 
                            sx={{ 
                              borderRadius: 1,
                              transition: 'all 0.2s',
                              '&:hover': {
                                backgroundColor: 'rgba(25, 118, 210, 0.04)'
                              }
                            }}
                          >
                            <Avatar sx={{ mr: 2, bgcolor: 'primary.light' }}>
                              {getFileIcon(filename)}
                            </Avatar>
                            <ListItemText 
                              primary={
                                <Typography variant="body1" fontWeight="medium">
                                  {filename}
                                </Typography>
                              }
                            />
                            <ListItemSecondaryAction>
                              <Tooltip title="Скачать файл">
                                <IconButton 
                                  onClick={() => download(filename)} 
                                  sx={{ 
                                    color: 'primary.main',
                                    transition: 'all 0.2s',
                                    '&:hover': {
                                      transform: 'translateY(-2px)',
                                      color: 'primary.dark'
                                    }
                                  }}
                                >
                                  <Download />
                                </IconButton>
                              </Tooltip>
                              <Tooltip title="Удалить файл">
                                <IconButton 
                                  onClick={() => confirmDelete(filename)}
                                  sx={{ 
                                    color: 'error.main',
                                    transition: 'all 0.2s',
                                    '&:hover': {
                                      transform: 'translateY(-2px)',
                                      color: 'error.dark'
                                    }
                                  }}
                                >
                                  <Delete />
                                </IconButton>
                              </Tooltip>
                            </ListItemSecondaryAction>
                          </ListItem>
                          {index < fileList.length - 1 && <Divider variant="inset" component="li" />}
                        </Box>
                      </Fade>
                    ))}
                  </List>
                )}
              </Paper>
            </Zoom>
          </Grid>
        </Grid>
        
        {/* Диалог подтверждения удаления */}
        <Dialog
          open={deleteDialogOpen}
          onClose={handleDeleteCancel}
          aria-labelledby="alert-dialog-title"
          aria-describedby="alert-dialog-description"
        >
          <DialogTitle id="alert-dialog-title">Подтверждение удаления</DialogTitle>
          <DialogContent>
            <DialogContentText id="alert-dialog-description">
              Вы действительно хотите удалить файл "{fileToDelete}"?
            </DialogContentText>
          </DialogContent>
          <DialogActions>
            <Button onClick={handleDeleteCancel} color="primary">
              Отмена
            </Button>
            <Button onClick={handleDeleteConfirm} color="error" autoFocus>
              Удалить
            </Button>
          </DialogActions>
        </Dialog>
        
        {/* Уведомления */}
        <Snackbar 
          open={!!error} 
          autoHideDuration={6000} 
          onClose={handleCloseSnackbar}
          anchorOrigin={{ vertical: 'bottom', horizontal: 'center' }}
          TransitionComponent={Slide}
        >
          <Alert 
            onClose={handleCloseSnackbar} 
            severity="error" 
            sx={{ width: '100%', boxShadow: 3 }}
            variant="filled"
          >
            {error}
          </Alert>
        </Snackbar>
        
        <Snackbar 
          open={!!success} 
          autoHideDuration={6000} 
          onClose={handleCloseSnackbar}
          anchorOrigin={{ vertical: 'bottom', horizontal: 'center' }}
          TransitionComponent={Slide}
        >
          <Alert 
            onClose={handleCloseSnackbar} 
            severity="success" 
            sx={{ width: '100%', boxShadow: 3 }}
            variant="filled"
          >
            {success}
          </Alert>
        </Snackbar>
      </Container>
    </Box>
  );
}