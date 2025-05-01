async function loadGoogleCharts(userMark, group, modal) {
  google.charts.load("current", { packages: ["corechart"] });

  let markCount = {};
  if (group === "science") {
    let lsValue;
    try {
      lsValue = await getData("09052024P2_scienceRankCount", "rank");
    } catch (e) {
      console.log("Error retrieving data from IndexedDB:", e);
      lsValue = null; // Set to null if there's an error}
    }
    //  = await getData("09052024P2_scienceRankCount", "rank");
    if (lsValue) {
      markCount = lsValue;
      console.log("Using cached data:", markCount);
    } else {
      await fetch(`/static/09052024P2/data/rank/scienceRankCount.json`)
        .then((response) => response.json())
        .then((data) => {
          markCount = data;
          storeData("09052024P2_scienceRankCount", "rank", data);
          console.log("Fetched and cached data:", markCount);
        });
    }
    // await fetch(`/static/09052024P2/data/rank/scienceRankCount.json`)
    //   .then((response) => response.json())
    //   .then((data) => {
    //     markCount = data;
    //   });
  } else {
    const lsValue = localStorage.getItem("09052024P2_" + "markMappingCount");
    if (lsValue) {
      markCount = JSON.parse(lsValue);
      console.log("Using cached data:", markCount);
    } else {
      await fetch(`/static/09052024P2/data/rank/markMappingCount.json`)
        .then((response) => response.json())
        .then((data) => {
          markCount = data;
          localStorage.setItem(
            "09052024P2_" + "markMappingCount",
            JSON.stringify(data)
          );
          console.log("Fetched and cached data:", markCount);
        });
    }

    // await fetch(`/static/09052024P2/data/rank/markMappingCount.json`)
    //   .then((response) => response.json())
    //   .then((data) => {
    //     markCount = data;
    //     console.log("Mark count data:", markCount);
    //   });
  }
  if (modal) {
    await drawChart();
  } else {
    google.charts.setOnLoadCallback(drawChart);
  }

  async function drawChart() {
    console.log("Drawing chart with userMark:", userMark);
    const data = new google.visualization.DataTable();
    data.addColumn("number", "Mark");
    data.addColumn("number", "Number of Students");
    data.addColumn("number", "Number of Students");

    const chartData = Object.entries(markCount).map(([mark, count]) => {
      const m = parseInt(mark);

      return [
        m,
        m === parseInt(userMark) ? null : count, // Default line skips 1000
        m === parseInt(userMark) ? count : null, // Only 1000 goes in highlight line
      ];
    });

    chartData.sort((a, b) => a[0] - b[0]);
    data.addRows(chartData);

    const options = {
      backgroundColor: "#2e2e2e",
      chartArea: {
        backgroundColor: "#2e2e2e",
        left: 100, // Enough space for Y-axis labels
        // right: 20,
        // top: 20,
        // bottom: 40,
        width: "80%",
        height: "70%",
      },
      legend: { position: "none" },
      colors: ["#ff8080", "#00ff7f"], // red for <1000, green for >1000
      hAxis: {
        title: "Mark",
        direction: -1,
        textStyle: { color: "#ffffff" },
        titleTextStyle: { color: "#ffffff" },
        gridlines: { color: "#444444" },
        viewWindow: {
          min: -1,
          max: 1201,
        },
      },
      vAxis: {
        title: "Count of Students",
        direction: 1,
        textStyle: { color: "#ffffff" },
        titleTextStyle: { color: "#ffffff" },
        gridlines: { color: "#444444" },
      },
    };
    if (modal) {
      // document.getElementById("fullscreen_chart").style.width = 2400 + "px";
      const chart = new google.visualization.ColumnChart(
        document.getElementById("fullscreen_chart")
      );
      options.width = 2400;
      options.height = 500;

      console.log("Drawing chart with data:", data);
      chart.draw(data, options);
    } else {
      const chart = new google.visualization.ColumnChart(
        document.getElementById(
          group === "science" ? "chart_div_science" : "chart_div"
        )
      );
      console.log("Drawing chart with data:", data);
      chart.draw(data, options);
    }
  }
}

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
        console.log("Transaction error:", event.target.error);
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
          console.log("Put error:", event.target.error);
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
      try {
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
