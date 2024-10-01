async function fetchAPI() {
  try {
    const res = await fetch(`${ENDPOINT}/${DATA_PATH}`, {
      headers: new Headers({
        Authorization: localStorage.getItem(LSATK),
      }),
    });
    const data = await res.json();
    Plotly.newPlot("plot", data);
    $("#loading").toggle();
  } catch (err) {
    console.error("Error fetching data:", err);
  }
}

$(document).ready(async () => await fetchAPI());
