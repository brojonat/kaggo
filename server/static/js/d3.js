import * as Plot from "https://cdn.jsdelivr.net/npm/@observablehq/plot@0.6/+esm";

async function fetchAPI() {
  // Read the data
  let karma_metrics = await d3.json(
    `${ENDPOINT}/timeseries/bucketed?` +
      new URLSearchParams({
        request_kind: "reddit.user",
        id: ID,
      }).toString(),
    {
      headers: new Headers({
        Authorization: localStorage.getItem(LSATK),
      }),
    }
  );
  let post_metrics = await d3.json(
    `${ENDPOINT}/metadata/children?` +
      new URLSearchParams({
        request_kind: "reddit.user",
        id: ID,
      }).toString(),
    {
      headers: new Headers({
        Authorization: localStorage.getItem(LSATK),
      }),
    }
  );
  karma_metrics = karma_metrics
    .filter((d) => {
      return d["metric"] === "reddit.user.link-karma";
    })
    .map((d) => {
      d["bucket"] = d3.isoParse(d["bucket"]);
      return d;
    });

  post_metrics = post_metrics.map((d) => {
    return {
      ts_created: d3.isoParse(d["data"]["ts_created"]),
      title: d["data"]["title"],
      link: d["data"]["link"],
    };
  });

  const options = {
    marks: [
      Plot.frame(),
      Plot.lineY(karma_metrics, {
        x: "bucket",
        y: "value",
        stroke: "metric",
        strokeWidth: 8,
        clip: true,
      }),
      Plot.ruleX(post_metrics, {
        x: "ts_created",
        clip: true,
      }),
      Plot.tip(
        post_metrics,
        Plot.pointer({
          x: "ts_created",
          title: (d) => d.link,
          lineWidth: 1000,
        })
      ),
    ],
    grid: true,
    inset: 10,
    facet: { marginRight: 90 },
    x: {
      tickSpacing: 80,
      label: "Time",
      domain: [new Date("2024-09-12"), new Date()],
    },
    y: {
      tickSpacing: 80,
      label: "Total Karma",
      tickFormat: d3.format("~s"),
    },
    // innerWidth: 500,
    // innerHeight: 500,
    legend: false,
    title: `Karma and Posts for ${ID}`,
  };
  const plot = Plot.plot(options);
  const div = document.querySelector("#plot");
  div.append(plot);
  $("#loading").toggle();
}

$(document).ready(async () => await fetchAPI());
