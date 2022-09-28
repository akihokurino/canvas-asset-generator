import * as functions from "firebase-functions";
const { CloudTasksClient } = require("@google-cloud/tasks");

const projectId = "canvas-329810";
const location = "asia-northeast1";

exports.queuingVideoSplitTask = functions
  .region(location)
  .storage.bucket("canvas-329810-video")
  .object()
  .onFinalize(async (object: functions.storage.ObjectMetadata) => {
    const filePath = object.name || "";
    const fileType = filePath.split(".").pop();
    if (fileType != "mp4") {
      console.log("mp4以外の拡張子はサポートしていません");
      return;
    }

    console.log("動画を分割するタスクをキューイングします");
    console.log(`ファイルパス ${filePath}`);

    const tasksClient = new CloudTasksClient();
    const queuePath = tasksClient.queuePath(projectId, location, "split-video");
    const url = "https://canvas-329810.an.r.appspot.com/split-video";
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
          Authorization: functions.config().token.internal,
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
