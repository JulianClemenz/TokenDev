
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

// --- Lógica de renderizado ---

/**
 * Renderiza la lista de ejercicios en el cuerpo de la tabla (Versión de Usuario)
 * @param {Array} exercises - La lista de ejercicios a mostrar.
 * @param {HTMLElement} tableBody - El elemento <tbody> de la tabla.
 */
function renderUserExercises(exercises, tableBody) {
  tableBody.innerHTML = '';

  if (exercises && exercises.length > 0) {
    exercises.forEach(exercise => {
      const exerciseId = exercise.id;

      const row = document.createElement('tr');

      row.innerHTML = `
        <td>${exercise.Name || ''}</td>
        <td>${exercise.Description || ''}</td>
        <td>${exercise.Category || ''}</td>
        <td>${exercise.MainMuscleGroup || ''}</td>
        <td>${exercise.DifficultLevel || ''}</td>
        <td><a href="${exercise.Example || '#'}" target="_blank" rel="noopener">Ver Video</a></td>
      `;
      tableBody.appendChild(row);
    });
  } else {
    tableBody.innerHTML = '<tr><td colspan="6">No se encontraron ejercicios con esos filtros.</td></tr>';
  }
}

// --- Lógica de la Página ---

/**
 * Carga los ejercicios desde la API (con o sin filtros) y los muestra en la tabla.
 */
async function loadExercises() {
  const tableBody = document.querySelector('.table tbody');
  tableBody.innerHTML = '<tr><td colspan="6">Cargando ejercicios...</td></tr>';

  // Leer valores de los filtros
  const name = document.getElementById('filter_name').value.trim();
  const category = document.getElementById('filter_category').value;
  const muscleGroup = document.getElementById('filter_muscle_group').value.trim();

  let endpoint = '';
  const params = new URLSearchParams();

  // Construir la URL del endpoint
  if (name) params.append('name', name);
  if (category) params.append('category', category);
  if (muscleGroup) params.append('muscle_group', muscleGroup);

  const queryString = params.toString();

  if (queryString) {
    endpoint = `/api/exercises/filter?${queryString}`;
  } else {
    endpoint = '/api/exercises';
  }

  try {
    const response = await fetchApi(endpoint);

    if (!response.ok) {
      const errorData = await response.json();
      throw new Error(errorData.error || `Error ${response.status}: No se pudieron cargar los ejercicios.`);
    }

    const exercises = await response.json();
    renderUserExercises(exercises, tableBody); // Usar la función de renderizado DE USUARIO

  } catch (error) {
    console.error('Error al cargar ejercicios:', error);
    if (error.message.includes("al menos un filtro")) {
      tableBody.innerHTML = '<tr><td colspan="6">No hay ejercicios para mostrar. Limpie los filtros para ver todos.</td></tr>';
    } else {
      tableBody.innerHTML = `<tr><td colspan="6" class="text-danger">Error: ${error.message}</td></tr>`;
    }
  }
}

// --- Inicialización ---

/**
 * Se ejecuta cuando el contenido del DOM está completamente cargado.
 */
document.addEventListener('DOMContentLoaded', () => {
  // Cargar la lista inicial (sin filtros)
  loadExercises();

  document.getElementById('btn_filter').addEventListener('click', loadExercises);

  document.getElementById('btn_clear_filters').addEventListener('click', () => {
    document.getElementById('filter_name').value = '';
    document.getElementById('filter_category').value = '';
    document.getElementById('filter_muscle_group').value = '';
    loadExercises(); // Recargar la lista completa
  });

});