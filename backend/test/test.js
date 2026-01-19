const fetch = require('node-fetch'); // Install with npm install node-fetch@2


const url = 'http://expired.badssl.com';

/*

First, it makes a request to the endpoint /start-scan/ to start generating the domain information, which returns the process ID.
Then, that ID is passed to the endpoint GET /scan-status/, where intervals are made until the information is completely generated.
At that point, the information is saved in the database by making a request to the endpoint POST /create-domain-info.
*/

async function demoAPI() {
  console.log('Iniciando demo de API...');

  // 1. POST a /start-scan
  const startResponse = await fetch('http://localhost:8080/start-scan', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ domain: url })
  });

  if (!startResponse.ok) {
    console.error('Error en POST start-scan:', startResponse.statusText);
    return;
  }

  const startData = await startResponse.json();
  const scanID = startData.scanRequestID;
  console.log('ID de escaneo obtenida:', scanID);

  // 2. Polling con setInterval hasta que status == "complete"
  let statusResponse;
  const pollInterval = setInterval(async () => {
    statusResponse = await fetch(`http://localhost:8080/scan-status/${scanID}`);
    const statusData = await statusResponse.json();

    if (statusData.status === 'complete') {
      clearInterval(pollInterval);
      console.log('Escaneo completado!');
      handleComplete(statusData);
    } else {
      console.log('Aún no están listos los resultados, esperando... (status actual:', statusData.status, ')');
    }
  }, 15000); // Polling cada 5 segundos
}

async function handleComplete(data) {
  const filteredResult = data.filteredResult;

  if (!filteredResult) {
    console.error('No se encontró filteredResult en la respuesta');
    return;
  }

  console.log('FilteredResult obtenido:', filteredResult);

  // 3. POST to save the filtered result in the MongoDB 
  const saveResponse = await fetch('http://localhost:8080/create-domain-info', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(filteredResult)
  });

  if (!saveResponse.ok) {
    console.error('Error guardando en DB:', saveResponse.statusText);
    return;
  }

  const saveData = await saveResponse.json();
  console.log('Registro guardado en DB:', saveData);
}

demoAPI();


