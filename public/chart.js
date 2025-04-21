async function loadGoogleCharts(userMark, group, modal) {
  google.charts.load("current", { packages: ["corechart"] });

  let markCount = {};
  if (group === "science") {
    await fetch(`/static/09052024P2/data/rank/scienceRankCount.json`)
      .then((response) => response.json())
      .then((data) => {
        markCount = data;
      });
  } else {
    await fetch(`/static/09052024P2/data/rank/markMappingCount.json`)
      .then((response) => response.json())
      .then((data) => {
        markCount = data;
        console.log("Mark count data:", markCount);
      });
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
      backgroundColor: "#252525",
      chartArea: {
        backgroundColor: "#252525",
        width: "85%",
        height: "70%",
      },
      legend: { position: "none" },
      colors: ["#ff8080", "#00ff7f"], // default red, highlight green
      hAxis: {
        direction: -1,
        textStyle: { color: "#ffffff" },
        titleTextStyle: { color: "#ffffff" },
        gridlines: { color: "#444444" },
        viewWindow: {
          min: -1, // lower than your lowest mark (950)
          max: 1201, // higher than your highest mark (1050)
        },
      },
      vAxis: {
        textStyle: { color: "#ffffff" },
        titleTextStyle: { color: "#ffffff" },
        gridlines: { color: "#444444" },
      },
    };
    if (modal) {
      const chart = new google.visualization.ColumnChart(
        document.getElementById("fullscreen_chart")
      );
      document.getElementById("modal-body").style.width = 2400 + "px";

      console.log("Drawing chart with data:", chartData);
      chart.draw(data, options);
    } else {
      const chart = new google.visualization.ColumnChart(
        document.getElementById(
          group === "science" ? "chart_div_science" : "chart_div"
        )
      );
      console.log("Drawing chart with data:", chartData);
      chart.draw(data, options);
    }
  }
}
