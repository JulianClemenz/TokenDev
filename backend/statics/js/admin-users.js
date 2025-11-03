
/**
 * Obtiene el token de autenticación desde sessionStorage.
 */
function getToken() {
  return sessionStorage.getItem('access_token');
}

/**
 * Realiza una solicitud fetch autenticada.
 * Esta función es necesaria porque GetLogs es una ruta de admin protegida.
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
    window.location.href = '/login'; // Asumiendo que el login está en la raíz
    throw new Error('No autorizado');
  }

  return response;
}

// --- Lógica de la Página ---

/**
 * Carga los logs de usuarios desde la API y los muestra en la tabla.
 */
async function loadUsers() {
  const tableBody = document.querySelector('.table tbody');
  tableBody.innerHTML = '<tr><td colspan="10">Cargando usuarios...</td></tr>';

  try {
    // Este es el endpoint que definiste en main.go para GetLogs
    const response = await fetchApi('/api/admin/stats/users');

    if (!response.ok) {
      // Manejo de 204 No Content (si el handler lo devuelve)
      if (response.status === 204) {
        tableBody.innerHTML = '<tr><td colspan="10">No hay usuarios registrados.</td></tr>';
        return;
      }
      throw new Error(`Error ${response.status}: No se pudieron cargar los usuarios.`);
    }

    const data = await response.json();
    tableBody.innerHTML = '';

    if (data.users && data.users.length > 0) {
      data.users.forEach((user, index) => {
        const row = document.createElement('tr');

        // Formatear la fecha (asumiendo que llega como ISO string "YYYY-MM-DDTHH...")
        const birthDate = user.BirthDate ? user.BirthDate.split('T')[0] : 'N/D';

        row.innerHTML = `
          <th scope="row">${index + 1}</th>
          <td>${user.UserName || ''}</td>
          <td>${user.Name || ''}</td>
          <td>${user.LastName || ''}</td>
          <td>${user.Email || ''}</td>
          <td>${birthDate}</td>
          <td>${user.Height || 0} cm</td>
          <td>${user.Weight || 0} kg</td>
          <td>${user.Experience || ''}</td>
          <td>${user.Objetive || ''}</td>
        `;
        tableBody.appendChild(row);
      });
    } else {
      tableBody.innerHTML = '<tr><td colspan="10">No hay usuarios registrados.</td></tr>';
    }
  } catch (error) {
    console.error('Error al cargar usuarios:', error);
    tableBody.innerHTML = `<tr class="text-center"><td colspan="10" class="text-danger">Error: ${error.message}</td></tr>`;
  }
}


// --- Inicialización ---
document.addEventListener('DOMContentLoaded', () => {
  loadUsers();
});