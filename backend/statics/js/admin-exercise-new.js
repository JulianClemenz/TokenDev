// --- Funciones de Ayuda (reutilizadas de admin-exercises.js) ---

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

// --- Lógica de la Página de Nuevo Ejercicio ---

/**
 * Maneja el envío del formulario para crear un nuevo ejercicio.
 */
async function handleSaveExercise() {
  const errorElement = document.getElementById('error_msg');
  errorElement.textContent = ''; // Limpiar errores previos

  // 1. Recolectar los datos del formulario
  // Los nombres de la API (ej: main_muscle_group) deben coincidir con tu ExcerciseRegisterDTO
  const payload = {
    name: document.getElementById('ex_name').value.trim(),
    main_muscle_group: document.getElementById('ex_group').value.trim(), // ID 'ex_group' mapea a 'main_muscle_group'
    description: document.getElementById('ex_desc').value.trim(),
    category: document.getElementById('ex_category').value,
    difficult_level: document.getElementById('ex_difficulty').value,
    example: document.getElementById('ex_sample').value.trim(),
    instructions: document.getElementById('ex_instructions').value.trim(),
  };

  // 2. Validación simple (el backend hará la validación final)
  if (!payload.name || !payload.main_muscle_group || !payload.description || !payload.category || !payload.difficult_level) {
    errorElement.textContent = 'Error: Debes completar todos los campos (Nombre, Grupo muscular, Descripción, Categoría y Dificultad son obligatorios).';
    return;
  }

  try {
    // 3. Enviar al endpoint (POST /api/exercises)
    const response = await fetchApi('/api/exercises', {
      method: 'POST',
      body: JSON.stringify(payload),
    });

    if (!response.ok) {
      // Si el servidor responde con un error (ej: 400, 409)
      const errorData = await response.json();
      throw new Error(errorData.error || 'Ocurrió un error al guardar.');
    }

    // 4. Éxito
    alert('¡Ejercicio creado exitosamente!');
    window.location.href = 'admin-exercises.html'; // Redirigir de vuelta al listado

  } catch (error) {
    console.error('Error al crear ejercicio:', error);
    errorElement.textContent = error.message;
  }
}

// --- Inicialización ---
document.addEventListener('DOMContentLoaded', () => {
  // 1. Asignar el evento al botón de guardar
  const saveButton = document.getElementById('btn_save_exercise');
  if (saveButton) {
    saveButton.addEventListener('click', handleSaveExercise);
  }
});