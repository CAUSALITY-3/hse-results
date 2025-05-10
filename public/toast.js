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

document.body.addEventListener("htmx:responseError", function (evt) {
  if (evt.detail.xhr.status === 429) {
    showRateLimitToast();
  }
});
