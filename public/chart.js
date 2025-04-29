async function loadGoogleCharts(userMark, group, modal) {
  google.charts.load("current", { packages: ["corechart"] });

  let markCount = {};
  if (group === "science") {
    const lsValue = localStorage.getItem("09052024P2_" + "scienceRankCount");
    if (lsValue) {
      markCount = JSON.parse(lsValue);
      console.log("Using cached data:", markCount);
    } else {
      await fetch(`/static/09052024P2/data/rank/scienceRankCount.json`)
        .then((response) => response.json())
        .then((data) => {
          markCount = data;
          localStorage.setItem(
            "09052024P2_" + "scienceRankCount",
            JSON.stringify(data)
          );
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
