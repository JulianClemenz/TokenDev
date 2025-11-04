// --- Variables Globales ---
let userRoutines = []; // Almacenará las rutinas del usuario
let addExerciseModalInstance; // Almacenará la instancia del modal de Bootstrap

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
    logout(); // Redirigir al login si no hay usuario
    return null;
  }
  return JSON.parse(userStr);
}

/**
 * Función de logout (necesaria para getCurrentUser)
 */
function logout() {
  sessionStorage.removeItem('access_token');
  sessionStorage.removeItem('refresh_token');
  sessionStorage.removeItem('user');
  window.location.href = '/index.html';
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
      const row = document.createElement('tr');

      row.innerHTML = `
        <td>${exercise.Name || ''}</td>
        <td>${exercise.Description || ''}</td>
        <td>${exercise.Category || ''}</td>
        <td>${exercise.MainMuscleGroup || ''}</td>
        <td>${exercise.DifficultLevel || ''}</td>
        <td><a href="${exercise.Example || '#'}" target="_blank" rel="noopener">Ver Video</a></td>
        <td>
          <button 
            class="btn btn-outline-primary btn-sm btn-add-exercise" 
            data-exercise-id="${exercise.id}" 
            data-exercise-name="${exercise.Name}"
            data-bs-toggle="modal" 
            data-bs-target="#addExerciseModal">
            + Añadir
          </button>
        </td>
      `;
      tableBody.appendChild(row);
    });
  } else {
    tableBody.innerHTML = '<tr><td colspan="7">No se encontraron ejercicios con esos filtros.</td></tr>';
  }
}

// --- Lógica de Carga de Datos ---

/**
 * Carga las rutinas (solo las del usuario actual) para el dropdown del modal.
 */
async function loadUserRoutines() {
  const selectElement = document.getElementById('modal_routine_select');
  const currentUser = getCurrentUser();
  if (!currentUser) return;

  try {
    const response = await fetchApi('/api/routines');
    if (!response.ok) {
      if (response.status === 404) { // 404 es esperado si no tiene rutinas
        selectElement.innerHTML = '<option value="">No tienes rutinas creadas</option>';
        return;
      }
      throw new Error('No se pudieron cargar tus rutinas');
    }

    const allRoutines = await response.json();

    // Filtramos SÓLO las rutinas del usuario actual
    userRoutines = allRoutines.filter(r => r.CreatorUserID === currentUser.id);

    if (userRoutines.length > 0) {
      selectElement.innerHTML = ''; // Limpiar "cargando"
      userRoutines.forEach(routine => {
        const option = document.createElement('option');
        option.value = routine.ID;
        option.textContent = routine.Name;
        selectElement.appendChild(option);
      });
    } else {
      selectElement.innerHTML = '<option value="">No tienes rutinas creadas</option>';
    }

  } catch (error) {
    console.error('Error al cargar rutinas:', error);
    selectElement.innerHTML = `<option value="">Error al cargar</option>`;
  }
}

/**
 * Carga los ejercicios desde la API (con o sin filtros) y los muestra en la tabla.
 */
async function loadExercises() {
  const tableBody = document.querySelector('.table tbody');
  tableBody.innerHTML = '<tr><td colspan="7">Cargando ejercicios...</td></tr>';

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
      tableBody.innerHTML = '<tr><td colspan="7">No hay ejercicios para mostrar. Limpie los filtros para ver todos.</td></tr>';
    } else {
      tableBody.innerHTML = `<tr colspan="7" class="text-danger">Error: ${error.message}</td></tr>`;
    }
  }
}

/**
 * Maneja el guardado del ejercicio en la rutina desde el modal.
 */
async function handleSaveToRoutine() {
  const errorElement = document.getElementById('modal_error_msg');
  errorElement.textContent = '';

  // 1. Obtener datos del modal
  const exerciseId = document.getElementById('modal_exercise_id').value;
  const routineId = document.getElementById('modal_routine_select').value;
  const series = parseInt(document.getElementById('modal_series').value, 10);
  const reps = parseInt(document.getElementById('modal_reps').value, 10);
  const weight = parseFloat(document.getElementById('modal_weight').value) || 0;

  // 2. Validar
  if (!routineId) {
    errorElement.textContent = 'Debes seleccionar una rutina (o crear una si no tienes).';
    return;
  }
  if (!exerciseId) {
    errorElement.textContent = 'Error: No se seleccionó un ejercicio. Cierra el modal e inténtalo de nuevo.';
    return;
  }
  if (isNaN(series) || series <= 0 || isNaN(reps) || reps <= 0) {
    errorElement.textContent = 'Las series y repeticiones deben ser números mayores a 0.';
    return;
  }
  if (isNaN(weight) || weight < 0) {
    errorElement.textContent = 'El peso debe ser un número igual o mayor a 0.';
    return;
  }

  // 3. Preparar payload (debe coincidir con ExcerciseInRoutineDTO de Go)
  const payload = {
    exercise_id: exerciseId,
    repetitions: reps,
    series: series,
    weight: weight
  };

  try {
    // 4. Llamar a la API (POST /api/routines/:id/exercises)
    const response = await fetchApi(`/api/routines/${routineId}/exercises`, {
      method: 'POST',
      body: JSON.stringify(payload)
    });

    if (!response.ok) {
      const errorData = await response.json();
      throw new Error(errorData.error || 'No se pudo añadir el ejercicio.');
    }

    // 5. Éxito
    alert('¡Ejercicio añadido a la rutina exitosamente!');
    addExerciseModalInstance.hide(); // Ocultar el modal

    // Limpiar campos del modal para la próxima vez
    document.getElementById('modal_series').value = '3';
    document.getElementById('modal_reps').value = '10';
    document.getElementById('modal_weight').value = '0';
    errorElement.textContent = '';


  } catch (error) {
    console.error('Error al guardar en rutina:', error);
    errorElement.textContent = error.message;
  }
}


// --- Inicialización ---

/**
 * Se ejecuta cuando el contenido del DOM está completamente cargado.
 */
document.addEventListener('DOMContentLoaded', () => {
  // Obtener la instancia del modal de Bootstrap
  const modalElement = document.getElementById('addExerciseModal');
  if (modalElement) {
    addExerciseModalInstance = new bootstrap.Modal(modalElement);

    // Evento que se dispara ANTES de que el modal se muestre
    modalElement.addEventListener('show.bs.modal', (event) => {
      // 'event.relatedTarget' es el botón que abrió el modal
      const button = event.relatedTarget;

      const exerciseId = button.dataset.exerciseId;
      const exerciseName = button.dataset.exerciseName;

      //Poner los datos en el modal
      document.getElementById('modal_exercise_id').value = exerciseId;
      document.getElementById('modal_exercise_name').textContent = exerciseName;

      // Limpiar errores previos
      document.getElementById('modal_error_msg').textContent = '';
    });
  }

  // Cargar la lista inicial (sin filtros)
  loadExercises();
  // Cargar las rutinas del usuario para el dropdown
  loadUserRoutines();

  // --- Asignar eventos a botones ---
  document.getElementById('btn_filter').addEventListener('click', loadExercises);

  document.getElementById('btn_clear_filters').addEventListener('click', () => {
    document.getElementById('filter_name').value = '';
    document.getElementById('filter_category').value = '';
    document.getElementById('filter_muscle_group').value = '';
    loadExercises(); // Recargar la lista completa
  });

  // Asignar evento al botón de GUARDAR del modal
  document.getElementById('btn_save_to_routine').addEventListener('click', handleSaveToRoutine);

});