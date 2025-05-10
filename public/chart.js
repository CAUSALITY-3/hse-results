async function loadGoogleCharts(userMark, mainRoute, group, elementId) {
  google.charts.load("current", { packages: ["corechart"] });

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
      await fetch(`/static/${mainRoute}/data/rank/markMappingCount.json`)
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
        `/static/${mainRoute}/data/schoolRank/${schoolCode}_rankCount.json`
      )
        .then((response) => response.json())
        .then(async (data) => {
          markCount = data;
          await storeData(key, "main", data);
          console.log("Fetched and cached data:", markCount);
        });
    }
    setRankDetails(userMark, "school", markCount);
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
      await fetch(`/static/${mainRoute}/data/rank/${groupValue}RankCount.json`)
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
      width: 1200,
      height: 500,
      legend: { position: "none" },
      colors: ["#ff8080", "#00ff7f"], // red for <1000, green for >1000
      hAxis: {
        title: "Mark",
        direction: -1,
        textStyle: { color: "#ababab" },
        titleTextStyle: { color: "#ffffff" },
        gridlines: { color: "#444444" },
        viewWindow: {
          min: elementId === "chart_div" ? 1 : -2,
          max: 1200,
        },
      },
      vAxis: {
        title: "Count of Students",
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

  function setRankDetails(userMark, group, markCount) {
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

    if (group === "overall") {
      console.log("Overall group rank:", rank);
    }
    const rankElement = document.getElementById(`student-${group}-rank`);
    if (rankElement) rankElement.innerText = rank;

    const totalCountElement = document.getElementById(
      `student-total-${group}-count`
    );
    if (totalCountElement) totalCountElement.innerText = totalCount;

    const aheadCountElement = document.getElementById(
      `student-total-${group}-ahead-count`
    );
    if (aheadCountElement) aheadCountElement.innerText = rank;

    const behindCountElement = document.getElementById(
      `student-total-${group}-behind-count`
    );
    if (behindCountElement) behindCountElement.innerText = behindMe;

    const sameCountElement = document.getElementById(
      `student-total-${group}-same-count`
    );
    if (sameCountElement) sameCountElement.innerText = rankCount;
  }
}
