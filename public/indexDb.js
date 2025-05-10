function storeData(key, collection, value) {
  const indexedDB =
    window.indexedDB ||
    window.mozIndexedDB ||
    window.webkitIndexedDB ||
    window.msIndexedDB ||
    window.shimIndexedDB;

  if (!indexedDB) {
    console.log("IndexedDB could not be found in this browser.");
  }
  const version = +localStorage.getItem("dbVersion") || 1;
  return new Promise((resolve, reject) => {
    const request = indexedDB.open("hseResults", version);

    request.onupgradeneeded = (event) => {
      const db = event.target.result;

      if (!db.objectStoreNames.contains(collection)) {
        console.log("Creating object store:", collection);
        db.createObjectStore(collection, { keyPath: "id" });
      }
    };

    request.onsuccess = (event) => {
      document.getElementById("school-search-route").innerText =
        "Rajan " + version;
      const db = event.target.result;

      console.log("collections", db.objectStoreNames.length);
      if (!db.objectStoreNames.contains(collection)) {
        console.log(`Object store "${collection}" not found.`);
        db.close();
        localStorage.setItem("dbVersion", version + 1);

        return storeData(key, collection, value, version + 1);
      }
      const tx = db.transaction(collection, "readwrite");
      tx.onerror = function (event) {
        db.close();
        console.log("Transaction error:", event.target.error);
      };

      const store = tx.objectStore(collection);
      try {
        const addRequest = store.put({ id: key, value: value });
        addRequest.onsuccess = function (event) {
          console.log("Put success:", event.target.result);
          resolve("Data stored successfully.");
        };
        addRequest.onerror = function (event) {
          console.log("Put error:", event.target.error);
          db.close();
          reject("Failed to write to the database.", event.target.error);
        };
      } catch (e) {
        db.close();
        if (e.name === "QuotaExceededError") {
          reject("Not enough storage space to save your data.");
        } else {
          reject("Unexpected error:", e);
        }
      }

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
      db.close();
      reject(event.target.error);
    };
    request.onblocked = (event) => {
      db.close();
      // Close all existing connections
      if (event.target.result) {
        event.target.result.close();
      }
    };
  });
}

function getData(key, collection) {
  const indexedDB =
    window.indexedDB ||
    window.mozIndexedDB ||
    window.webkitIndexedDB ||
    window.msIndexedDB ||
    window.shimIndexedDB;

  if (!indexedDB) {
    console.log("IndexedDB could not be found in this browser.");
  }
  const version = localStorage.getItem("dbVersion") || 1;
  return new Promise((resolve, reject) => {
    const request = indexedDB.open("hseResults", version);
    request.onupgradeneeded = (event) => {
      const db = event.target.result;
      if (!db.objectStoreNames.contains(collection)) {
        console.log("Creating object store:", collection);
        db.createObjectStore(collection, { keyPath: "id" });
      }
    };

    request.onsuccess = (event) => {
      try {
        const db = event.target.result;
        const tx = db.transaction(collection, "readonly");
        const store = tx.objectStore(collection);

        const getRequest = store.get(key);

        getRequest.onsuccess = () => {
          resolve(getRequest?.result?.value);
          console.log("Data retrieved successfully:", getRequest.result.value);
          db.close();
        };

        getRequest.onerror = () => {
          reject(getRequest.error);
          db.close();
        };
      } catch (e) {
        if (e.name === "NotFoundError") {
          console.log("Data not found in the database:", e);
          reject("Data not found in the database.");
        } else {
          console.log("Unexpected error:", e);
          reject("Unexpected error:", e);
        }
      }
    };

    request.onerror = (event) => {
      reject(event.target.error);
    };
  });
}
function deleteData(key, collection) {
  const indexedDB =
    window.indexedDB ||
    window.mozIndexedDB ||
    window.webkitIndexedDB ||
    window.msIndexedDB ||
    window.shimIndexedDB;

  if (!indexedDB) {
    console.log("IndexedDB could not be found in this browser.");
  }
  const version = localStorage.getItem("dbVersion") || 1;
  return new Promise((resolve, reject) => {
    const request = indexedDB.open("hse-results", version);

    request.onsuccess = (event) => {
      const db = event.target.result;
      const tx = db.transaction(collection, "readwrite");
      const store = tx.objectStore(collection);

      const deleteRequest = store.delete(key);

      deleteRequest.onsuccess = () => {
        resolve("Data deleted successfully.");
        db.close();
      };

      deleteRequest.onerror = () => {
        reject(deleteRequest.error);
      };
    };

    request.onerror = (event) => {
      reject(event.target.error);
      db.close();
    };
  });
}
