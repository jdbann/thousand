import "@hotwired/turbo";

import { Application } from "@hotwired/stimulus";

import FrameController from "./controllers/frame-controller.js";

window.Stimulus = Application.start();
Stimulus.register("frame", FrameController);
