// --- Funciones de Ayuda ---

/**
 * Obtiene el token de autenticación desde sessionStorage.
 */
function getToken() {
  return sessionStorage.getItem('access_token');
}

/**
 * Obtiene los datos del usuario (incluyendo el ID) desde sessionStorage.
 */
function getCurrentUser() {
    const userStr = sessionStorage.getItem('user');
    if (!userStr) {
        // Si no hay usuario, es un problema de autenticación
        logout(); // Redirigir al login
        return null;
    }
    return JSON.parse(userStr);
}


/**
 * Realiza una solicitud fetch autenticada.
 */
async function fetchApi(url, options = {}) {
  const token = getToken();
  
  const headers = {
    'Content-Type': 'application/json',
    'Authorization': `Bearer ${token}`,
    ...options.headers,
  };

  const response = await fetch(url, { ...options, headers });

  if (response.status === 401) {
    alert('Tu sesión ha expirado. Por favor, inicia sesión de nuevo.');
    window.location.href = '/login'; 
    throw new Error('No autorizado');
  }
  
  return response;
}

/**
 * Función de logout (usada por el helper si falla la autenticación)
 */
function logout() {
  sessionStorage.removeItem('access_token');
  sessionStorage.removeItem('refresh_token');
  sessionStorage.removeItem('user');
  window.location.href = '/index.html';
}


// --- Lógica de la Página ---

/**
 * Carga las rutinas del usuario en la tabla.
 */
async function loadRoutines() {
  const tableBody = document.getElementById('routines-table-body');
  const errorElement = document.getElementById('error_msg'); // Asegúrate de tener <p id="error_msg"></p> en tu HTML
  const currentUser = getCurrentUser();

  if (!currentUser || !currentUser.ID) {
    if (errorElement) errorElement.textContent = 'No se pudo identificar al usuario. Por favor, inicia sesión de nuevo.';
    // Aún si no podemos identificar al usuario, no deberíamos bloquear la carga
    // pero el filtro de abajo simplemente resultará en una lista vacía.
    console.warn("No se pudo obtener el ID del usuario desde sessionStorage.");
  }
  
  tableBody.innerHTML = '<tr><td colspan="4">Cargando tus rutinas...</td></tr>';

  try {
    const response = await fetchApi('/api/routines');
    if (!response.ok) {
        // Manejo de respuesta 204 (No Content) o 404 (Not Found)
        if (response.status === 204 || response.status === 404) {
            tableBody.innerHTML = '<tr><td colspan="4">Aún no has creado ninguna rutina.</td></tr>';
            return;
        }
        const errData = await response.json();
        throw new Error(errData.error || 'No se pudieron cargar las rutinas');
    }
    
    const routines = await response.json();
    
    // Fallback por si la API devuelve null en lugar de []
    if (!routines) {
        tableBody.innerHTML = '<tr><td colspan="4">Aún no has creado ninguna rutina.</td></tr>';
        return;
    }

    // Filtramos para mostrar SÓLO las rutinas creadas por el usuario actual
    const userRoutines = routines.filter(r => r.CreatorUserID === currentUser.ID);
    
    tableBody.innerHTML = ''; // Limpiar "cargando"
    
    if (userRoutines && userRoutines.length > 0) {
      userRoutines.forEach(routine => {
        const row = document.createElement('tr');
        // Usar EditionDate si existe, sino CreationDate
        const dateStr = routine.EditionDate || routine.CreationDate || new Date().toISOString();
        const date = new Date(dateStr);
        const formattedDate = date.toLocaleDateString('es-ES', {
            day: '2-digit',
            month: '2-digit',
            year: 'numeric'
        });
        
        // (Usamos 'ExcerciseList' y 'CreatorUserID' como en tu DTO de Go)
        // (Usamos 'ID' como en tu DTO de Go, que se mapea desde '_id')
        row.innerHTML = `
          <td><strong>${routine.Name}</strong></td>
          <td>${routine.ExcerciseList ? routine.ExcerciseList.length : 0}</td>
          <td>${formattedDate}</td>
          <td class="d-flex gap-2">
            <a href="user-routine-view.html?id=${routine.ID}" class="btn btn-outline-info btn-sm">Ver</a>
            <a href="user-routine-edit.html?id=${routine.ID}" class="btn btn-outline-primary btn-sm">Editar</a> 
            <button type="button" class="btn btn-outline-danger btn-sm btn-delete-routine" data-id="${routine.ID}">
              Eliminar
            </button>
          </td>
        `;
        tableBody.appendChild(row);
      });
    } else {
      tableBody.innerHTML = '<tr><td colspan="4">Aún no has creado ninguna rutina.</td></tr>';
    }

  } catch (error) {
    console.error('Error al cargar rutinas:', error);
    if (errorElement) errorElement.textContent = `Error: ${error.message}`;
    tableBody.innerHTML = `<tr><td colspan="4" class="text-danger">Error al cargar las rutinas.</td></tr>`;
  }
}

/**
 * Maneja el clic en el botón de eliminar rutina.
 */
async function handleDeleteRoutine(routineId) {
  if (!confirm('¿Estás seguro de que deseas eliminar esta rutina? Esta acción no se puede deshacer.')) {
    return;
  }

  const errorElement = document.getElementById('error_msg');
  if (errorElement) errorElement.textContent = ''; // Limpiar errores previos

  try {
    const response = await fetchApi(`/api/routines/${routineId}`, {
      method: 'DELETE'
    });

    if (!response.ok) {
      const errorData = await response.json();
      throw new Error(errorData.error || 'Error al eliminar la rutina.');
    }

    alert('Rutina eliminada correctamente.');
    loadRoutines(); // Recargar la tabla
    
  } catch (error) {
    console.error('Error al eliminar rutina:', error);
    if (errorElement) errorElement.textContent = `Error: ${error.message}`;
  }
}


// --- Inicialización ---
document.addEventListener('DOMContentLoaded', () => {
  // 1. Cargar las rutinas al iniciar
  loadRoutines();

  // 2. Escuchar clics en la tabla para los botones de eliminar
  const tableBody = document.getElementById('routines-table-body');
  tableBody.addEventListener('click', (event) => {
    const deleteButton = event.target.closest('.btn-delete-routine');
    
    if (deleteButton) {
      const routineId = deleteButton.dataset.id;
      handleDeleteRoutine(routineId);
    }
  });
});