import { Elm } from "./Main.elm";

let userStr = localStorage.getItem("user");
let user = null;

try {
	user = JSON.parse(userStr);
} catch (e) {
	console.error("Failed to parse user from local storage!", e);
}

Elm.Main.init({
	node: document.getElementById("app"),
	flags: { user },
});
