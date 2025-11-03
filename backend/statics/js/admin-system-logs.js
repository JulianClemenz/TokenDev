// --- Funciones de Ayuda ---

/**
 * Obtiene el token de autenticación desde sessionStorage.
 */
function getToken() {
  return sessionStorage.getItem('access_token');
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
    // Token inválido o expirado, redirigir al login
    alert('Tu sesión ha expirado. Por favor, inicia sesión de nuevo.');
    window.location.href = '/login'; 
    throw new Error('No autorizado');
  }
  
  return response;
}

// --- Lógica de la Página ---

/**
 * Carga los usuarios ACTIVOS desde la API y los muestra en la tabla.
 */
async function loadActiveUsers() {
  const tableBody = document.querySelector('.table tbody');
  // Ajustamos el colspan a 4, según las nuevas columnas
  tableBody.innerHTML = '<tr><td colspan="4">Cargando usuarios activos...</td></tr>'; 

  try {
    // 1. Llamamos al mismo endpoint que usa admin-users
    const response = await fetchApi('/api/admin/stats/users'); 
    
    if (response.status === 204) {
      tableBody.innerHTML = '<tr><td colspan="4">No hay usuarios en el sistema.</td></tr>';
      return;
    }
    
    if (!response.ok) {
      throw new Error(`Error ${response.status}: No se pudieron cargar los usuarios.`);
    }

    const data = await response.json(); // La respuesta es { total: N, users: [...] }
    
    // 2. FILTRAMOS por IsActive = true
    const activeUsers = data.users.filter(user => user.is_active === true);

    tableBody.innerHTML = ''; // Limpiar "cargando"

    if (activeUsers.length > 0) {
      activeUsers.forEach(user => {
        const row = document.createElement('tr');
        
        // 3. APLICAMOS EL ESTILO VERDE (Clase de Bootstrap para "éxito")
        row.className = 'table-success'; 

        row.innerHTML = `
          <td>${user.UserName || 'N/D'}</td>
          <td>${user.Email || 'N/D'}</td>
          <td>${user.Role || 'N/D'}</td>
          <td><span class="badge bg-success">Activo</span></td>
        `;
        tableBody.appendChild(row);
      });
    } else {
      tableBody.innerHTML = '<tr><td colspan="4">No hay usuarios activos en este momento.</td></tr>';
    }
  } catch (error) {
    console.error('Error al cargar usuarios activos:', error);
    tableBody.innerHTML = `<tr class="text-center"><td colspan="4" class="text-danger">Error: ${error.message}</td></tr>`;
  }
}


// --- Inicialización ---

/**
 * Se ejecuta cuando el contenido del DOM está completamente cargado.
 */
document.addEventListener('DOMContentLoaded', () => {
  loadActiveUsers();
});