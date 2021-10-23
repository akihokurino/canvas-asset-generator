import * as functions from "firebase-functions";
import * as path from "path";
const { CloudTasksClient } = require("@google-cloud/tasks");

const projectId = "canvas-329810";
const location = "asia-northeast1";
const queue = "video-split-task";

exports.queuingVideoSplitTask = functions
  .region(location)
  .storage.object()
  .onFinalize(async (object: functions.storage.ObjectMetadata) => {
    const filePath = object.name || "";
    const fileDir = path.dirname(filePath);

    if (fileDir !== "Video") {
      return;
    }

    console.log("動画を分割するタスクをキューイングします");
    console.log(`ファイルパス ${filePath}`);
    console.log(`ファイルディレクトリ ${fileDir}`);

    const tasksClient = new CloudTasksClient();
    const queuePath = tasksClient.queuePath(projectId, location, queue);

    const url = "https://canvas-329810.an.r.appspot.com/video_split";
    const delaySeconds = 1;

    const payload = {
      path: filePath,
    };

    const task = {
      httpRequest: {
        httpMethod: "POST",
        url,
        body: Buffer.from(JSON.stringify(payload)).toString("base64"),
        headers: {
          "Content-Type": "application/json",
        },
      },
      scheduleTime: {
        seconds: delaySeconds,
      },
    };

    const [response] = await tasksClient.createTask({
      parent: queuePath,
      task,
    });

    console.log("キューイングが完了しました。");
    console.log(`task name = ${response.name}`);
  });
