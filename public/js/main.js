// Function to handle zoom adjustment
function adjustZoom(elementIds) {
  let pageWidth = 1400;
  let screenWidth = window.innerWidth;

  console.log(pageWidth, screenWidth, window);
  if (!elementIds) return;
  for (const elementId of elementIds) {
    const element = document.getElementById(elementId);
    if (!element) {
      console.error("Element with ID 'resultSection' not found.");
      continue;
    }
    if (screenWidth < pageWidth) {
      console.log(pageWidth, screenWidth, "kklklklklk");
      let scaleFactor = screenWidth / pageWidth;
      element.style.zoom = scaleFactor;
      console.log(scaleFactor);
      // document.body.style.transformOrigin = 'top left';
    } else {
      element.style.zoom = 1;
    }
  }
}

// Function to load Google Charts and draw the chart
async function loadGoogleCharts(userMark, mainRoute, group, elementId) {
  google.charts.load("current", { packages: ["corechart"] });

  console.log("assssddddddd", userMark, mainRoute, group, elementId);
  let markCount = {};
  if (group === "overall") {
    let lsValue;
    try {
      lsValue = await getData(`${mainRoute}_markMappingCount`, "main");
    } catch (e) {
      console.log("Error retrieving data from IndexedDB:", e);
      lsValue = null; // Set to null if there's an error
    }
    if (lsValue) {
      markCount = lsValue;
      console.log("Using cached data:", markCount);
    } else {
      await fetch(`/static/${mainRoute}/rank/markMappingCount.json`)
        .then((response) => response.json())
        .then(async (data) => {
          markCount = data;
          await storeData(`${mainRoute}_markMappingCount`, "main", data);
          console.log("Fetched and cached data:", markCount);
        });
    }
    setRankDetails(userMark, "overall", markCount);
  } else if (group.startsWith("school_")) {
    const schoolCode = group.split("_")[1];
    console.log("zxcvbbnn", mainRoute + "_" + schoolCode + "_rankCount");
    const key = mainRoute + "_" + schoolCode + "_rankCount";
    let lsValue;
    try {
      lsValue = await getData(key, "main");
    } catch (e) {
      console.log("Error retrieving data from IndexedDB:", e);
      lsValue = null; // Set to null if there's an error
    }
    if (lsValue) {
      markCount = lsValue;
      console.log("Using cached data:", markCount);
    } else {
      await fetch(
        `/static/${mainRoute}/schoolRank/${schoolCode}_rankCount.json`
      )
        .then((response) => response.json())
        .then(async (data) => {
          markCount = data;
          await storeData(key, "main", data);
          console.log("Fetched and cached data:", markCount);
        });
    }
    setRankDetails(userMark, "school", markCount);
  } else if (group === "passPercentageCount") {
    let lsValue;
    const key = mainRoute + "_" + group;
    try {
      lsValue = await getData(key, "main");
    } catch (e) {
      console.log("Error retrieving data from IndexedDB:", e);
      lsValue = null;
    }
    if (lsValue) {
      markCount = lsValue;
      console.log("Using cached data:", markCount);
    } else {
      await fetch(`/static/${mainRoute}/rank/${group}.json`)
        .then((response) => response.json())
        .then(async (data) => {
          markCount = data;
          await storeData(key, "main", data);
          console.log("Fetched and cached data:", markCount);
        });
    }
    setRankDetails(userMark, "overall", markCount, true);
  } else {
    let groupValue = group.toLowerCase();
    let lsValue;
    const key = mainRoute + "_" + groupValue + "RankCount";
    try {
      lsValue = await getData(key, "main");
    } catch (e) {
      console.log("Error retrieving data from IndexedDB:", e);
      lsValue = null;
    }
    if (lsValue) {
      markCount = lsValue;
      console.log("Using cached data:", markCount);
    } else {
      await fetch(`/static/${mainRoute}/rank/${groupValue}RankCount.json`)
        .then((response) => response.json())
        .then(async (data) => {
          markCount = data;
          await storeData(key, "main", data);
          console.log("Fetched and cached data:", markCount);
        });
    }
    setRankDetails(userMark, "group", markCount);
  }
  console.log(elementId, "markCount:");
  google.charts.setOnLoadCallback(() => drawChart(elementId));

  async function drawChart(elementId) {
    console.log("Drawing chart with userMark:", userMark, elementId);
    const data = new google.visualization.DataTable();
    data.addColumn("number", "Mark");
    data.addColumn(
      "number",
      elementId === "chart_div_school_page"
        ? "Number of Schools"
        : "Number of Students"
    );
    data.addColumn(
      "number",
      elementId === "chart_div_school_page"
        ? "Number of Schools"
        : "Number of Students"
    );

    const chartData = Object.entries(markCount).map(([mark, count]) => {
      const m = parseInt(mark);

      return !userMark
        ? [m, count, null]
        : [
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
        left: elementId === "chart_div_school_page" ? 100 : 100,
        right: elementId === "chart_div_school_page" ? -200 : 0,
        ...(elementId === "chart_div_school_page"
          ? {}
          : { width: "80%", height: "70%" }),
      },
      width: elementId === "chart_div_school_page" ? 1200 : 1200,
      height: 500,
      legend: { position: "none" },
      colors: ["#ff8080", "#00ff7f"], // red for <1000, green for >1000
      hAxis: {
        title:
          elementId === "chart_div_school_page" ? "Pass Percentage" : "Mark",
        direction: -1,
        textStyle: { color: "#ababab" },
        titleTextStyle: { color: "#ffffff" },
        gridlines: { color: "#444444" },
        viewWindow: {
          min: elementId === "chart_div" ? 1 : -2,
          max: elementId === "chart_div_school_page" ? 100 : 1200,
        },
      },
      vAxis: {
        title:
          elementId === "chart_div_school_page"
            ? "Count of Schools"
            : "Count of Students",
        direction: 1,
        textStyle: { color: "#ababab" },
        titleTextStyle: { color: "#ffffff" },
        gridlines: { color: "#444444" },
      },
    };
    const chart = new google.visualization.ColumnChart(
      document.getElementById(elementId)
    );
    console.log("Drawing chart with data:", data);
    chart.draw(data, options);
  }

  function setRankDetails(userMark, group, markCount, school = false) {
    const rankCount = markCount[userMark] || 0;
    const totalCount = Object.values(markCount).reduce((a, b) => a + b, 0);
    let rank = 0;
    let behindMe = 0;
    for (const mark in markCount) {
      if (parseInt(mark) > parseInt(userMark)) {
        rank += markCount[mark];
      } else if (parseInt(mark) < parseInt(userMark)) {
        behindMe += markCount[mark]; // Count the number of students behind the user
      }
    }

    let prefix = "student";
    if (school) {
      prefix = "school";
    }
    if (group === "overall") {
      console.log("Overall group rank:", rank);
    }
    console.log(prefix, group);
    const rankElement = document.getElementById(`${prefix}-${group}-rank`);
    console.log({ rankElement: `${prefix}-${group}-rank` });
    if (rankElement) rankElement.innerText = rank;

    const totalCountElement = document.getElementById(
      `${prefix}-total-${group}-count`
    );
    console.log({ totalCountElement: `${prefix}-${group}-rank` });
    if (totalCountElement) totalCountElement.innerText = totalCount;

    const aheadCountElement = document.getElementById(
      `${prefix}-total-${group}-ahead-count`
    );
    if (aheadCountElement) aheadCountElement.innerText = rank;

    const behindCountElement = document.getElementById(
      `${prefix}-total-${group}-behind-count`
    );
    if (behindCountElement) behindCountElement.innerText = behindMe;

    const sameCountElement = document.getElementById(
      `${prefix}-total-${group}-same-count`
    );
    if (sameCountElement) sameCountElement.innerText = rankCount;
  }
}

// Function to store data in IndexedDB
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
          console.log("Data retrieved successfully:", key);
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

// Function to show a toast message for rate limiting
function showRateLimitToast(
  message = "⚠️ Too many requests. Please wait a few seconds.",
  duration = 5000
) {
  // If toast or overlay already exists, remove them first
  const existingToast = document.getElementById("rate-limit-toast");
  const existingOverlay = document.getElementById("toast-overlay");
  if (existingToast) existingToast.remove();
  if (existingOverlay) existingOverlay.remove();

  // Create overlay element
  const overlay = document.createElement("div");
  overlay.id = "toast-overlay";

  // Style the overlay
  Object.assign(overlay.style, {
    position: "fixed",
    top: "0",
    left: "0",
    width: "100%",
    height: "100%",
    backgroundColor: "rgba(0, 0, 0, 0.5)", // Semi-transparent black
    zIndex: 999, // Below the toast
  });

  // Prevent clicks on the overlay
  overlay.addEventListener("click", (e) => e.stopPropagation());

  // Create toast element
  const toast = document.createElement("div");
  toast.id = "rate-limit-toast";

  // Create spinner element
  const spinner = document.createElement("div");
  spinner.className = "toast-spinner";

  // Style the spinner
  Object.assign(spinner.style, {
    width: "24px",
    height: "24px",
    border: "4px solid white",
    borderTop: "4px solid transparent",
    borderRadius: "50%",
    animation: "spin 1s linear infinite",
    marginRight: "12px",
  });

  // Add spinner animation
  const style = document.createElement("style");
  style.textContent = `
        @keyframes spin {
            0% { transform: rotate(0deg); }
            100% { transform: rotate(360deg); }
        }
    `;
  document.head.appendChild(style);

  // Create message container
  const messageContainer = document.createElement("span");
  messageContainer.innerText = message;

  // Style the toast
  Object.assign(toast.style, {
    display: "flex",
    alignItems: "center",
    position: "fixed",
    bottom: "20px",
    right: "20px",
    backgroundColor: "#f87171",
    color: "white",
    padding: "16px 24px",
    borderRadius: "8px",
    boxShadow: "0 4px 12px rgba(0,0,0,0.2)",
    fontFamily: "sans-serif",
    zIndex: 1000, // Above the overlay
    opacity: "0",
    transition: "opacity 0.3s ease",
  });

  // Append spinner and message to toast
  toast.appendChild(spinner);
  toast.appendChild(messageContainer);

  document.body.appendChild(overlay);
  document.body.appendChild(toast);

  // Fade in
  setTimeout(() => (toast.style.opacity = "1"), 10);

  // Fade out and remove
  setTimeout(() => {
    toast.style.opacity = "0";
    setTimeout(() => {
      toast.remove();
      overlay.remove();
    }, 1000);
  }, duration);
}

document?.body?.addEventListener("htmx:responseError", function (evt) {
  if (evt?.detail?.xhr?.status === 429) {
    showRateLimitToast();
  }
});
