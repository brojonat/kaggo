async function fetchAPI() {
  try {
    const res = await fetch(
      `${ENDPOINT}/plot-data?id=workflow-count-projection`,
      {
        headers: new Headers({
          Authorization: localStorage.getItem(LSATK),
        }),
      }
    );
    const data = await res.json();
    Plotly.newPlot("plot", data);
    $("#loading").toggle();
  } catch (err) {
    console.error("Error fetching data:", err);
  }
}

$(document).ready(async () => await fetchAPI());
