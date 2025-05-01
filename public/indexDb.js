const indexedDB =
  window.indexedDB ||
  window.mozIndexedDB ||
  window.webkitIndexedDB ||
  window.msIndexedDB ||
  window.shimIndexedDB;

if (!indexedDB) {
  console.log("IndexedDB could not be found in this browser.");
}

function storeData(key, collection, value) {
  return new Promise((resolve, reject) => {
    const request = indexedDB.open("hse-results", 1);

    request.onupgradeneeded = (event) => {
      const db = event.target.result;
      if (!db.objectStoreNames.contains(collection)) {
        db.createObjectStore(collection, { keyPath: "id" });
      }
    };

    request.onsuccess = (event) => {
      const db = event.target.result;
      if (!db.objectStoreNames.contains(collection)) {
        reject(`Object store "${collection}" not found.`);
        db.close();
        return;
      }
      const tx = db.transaction(collection, "readwrite");

      tx.onerror = function (event) {
        console.error("Transaction error:", event.target.error);
        alert("Storage operation failed. Possibly no space left.");
      };

      const store = tx.objectStore(collection);

      try {
        const addRequest = store.put({ id: key, value: value });
        addRequest.onsuccess = function (event) {
          console.log("Put success:", event.target.result);
          resolve("Data stored successfully.");
        };
        addRequest.onerror = function (event) {
          console.error("Put error:", event.target.error);
          reject("Failed to write to the database.", event.target.error);
        };
      } catch (e) {
        if (e.name === "QuotaExceededError") {
          reject("Not enough storage space to save your data.");
        } else {
          reject("Unexpected error:", e);
        }
      }

      //   store.put({ id: key, value: value }); // stores the JSON object directly
      //   resolve("Data stored successfully.");
      tx.oncomplete = () => {
        db.close();
        console.log("Transaction completed: database modification finished.");
      };
      tx.onerror = () => {
        db.close();
        console.log("Transaction not completed: database modification failed.");
      };
    };

    request.onerror = (event) => {
      reject(event.target.error);
    };
  });
}

function getData(key, collection) {
  return new Promise((resolve, reject) => {
    const request = indexedDB.open("hse-results", 1);

    request.onsuccess = (event) => {
      const db = event.target.result;
      const tx = db.transaction(collection, "readonly");
      const store = tx.objectStore(collection);

      const getRequest = store.get(key);

      getRequest.onsuccess = () => {
        resolve(getRequest.result.value);
        console.log("Data retrieved successfully:", getRequest.result.value);
        db.close();
      };

      getRequest.onerror = () => {
        reject(getRequest.error);
        db.close();
      };
    };

    request.onerror = (event) => {
      reject(event.target.error);
    };
  });
}
function deleteData(key, collection) {
  return new Promise((resolve, reject) => {
    const request = indexedDB.open("hse-results", 1);

    request.onsuccess = (event) => {
      const db = event.target.result;
      const tx = db.transaction(collection, "readwrite");
      const store = tx.objectStore(collection);

      const deleteRequest = store.delete(key);

      deleteRequest.onsuccess = () => {
        resolve("Data deleted successfully.");
      };

      deleteRequest.onerror = () => {
        reject(deleteRequest.error);
      };
    };

    request.onerror = (event) => {
      reject(event.target.error);
    };
  });
}
