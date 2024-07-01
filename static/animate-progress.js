function animateProgressBar() {
  // elements
  let progressBar = document.querySelector("progress");
  let completedCounter = document.querySelector("#completed-counter");
  let countElement;
  if (completedCounter != null) {
    countElement = completedCounter.children.item(1);
  }
  // original to animate to
  let original = parseFloat(progressBar.getAttribute("data-value"));
  // max to calculate completed persentage of current value
  let max = parseFloat(progressBar.getAttribute("max"));
  // old value to animate from
  let old = parseFloat(progressBar.getAttribute("data-old-completed"));
  // calculate position right based on multiple values
  let right = (value) => {
    counterWidth = 0.04 * progressBar.clientWidth;
    // convert from persentage to value in pixel using progressbar width
    return (parseFloat(value) / max) * progressBar.clientWidth - counterWidth;
  };
  // set value to old value to start animation
  // -- set current to old
  progressBar.value = old;
  if (completedCounter != null) {
    // show count label
    completedCounter.style.display = `flex`;
    // set starting right position to calculated right
    completedCounter.style.right = `${right(progressBar.value)}px`;
    countElement.innerText = Math.ceil(progressBar.value);
  }
  // animation timer
  let timer = 5;
  function updateProgress() {
    if (completedCounter != null) {
      completedCounter.style.right = `${right(progressBar.value)}px`;
      countElement.innerText = Number(progressBar.value).toFixed(0);
    }
    if (original < old) {
      if (progressBar.value > original) {
        progressBar.value -= 0.02;
        setTimeout(updateProgress, timer);
      }
    } else {
      if (progressBar.value < original) {
        progressBar.value += 0.02;
        setTimeout(updateProgress, timer);
      }
    }
  }
  updateProgress();
}

animateProgressBar();
