package browser

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func stopRecording(cli *client.Client, containerID string, sessionId int64) error {
	pkillCmd := []string{"pkill", "-SIGINT", "ffmpeg"}
	execResp, err := cli.ContainerExecCreate(context.Background(), containerID, container.ExecOptions{
		User:         "0:0",
		Cmd:          pkillCmd,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          false,
		Detach:       false,
	})
	if err != nil {
		log.Printf("Failed to create pkill exec for session %d: %v", sessionId, err)
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := cli.ContainerExecStart(ctx, execResp.ID, container.ExecStartOptions{}); err != nil {
		log.Printf("Failed to send SIGINT to ffmpeg: %v", err)
	}

	log.Printf("Sent SIGINT to ffmpeg process for session %d", sessionId)

	log.Printf("Waiting for ffmpeg to finish writing video file for session %d...", sessionId)

	// Check if ffmpeg is still running and wait for it to finish
	for i := 0; i < 10; i++ { // Wait up to 10 seconds
		checkCmd := []string{"pgrep", "ffmpeg"}
		checkExecResp, err := cli.ContainerExecCreate(context.Background(), containerID, container.ExecOptions{
			User:         "0:0",
			Cmd:          checkCmd,
			AttachStdout: true,
			AttachStderr: true,
			Tty:          false,
			Detach:       false,
		})
		if err != nil {
			log.Printf("Failed to check ffmpeg status: %v", err)
			break
		}

		if err := cli.ContainerExecStart(ctx, checkExecResp.ID, container.ExecStartOptions{}); err != nil {
			log.Printf("Failed to execute ffmpeg check: %v", err)
			break
		}

		checkInspectResp, err := cli.ContainerExecInspect(ctx, checkExecResp.ID)
		if err != nil {
			log.Printf("Failed to inspect ffmpeg check: %v", err)
			break
		}

		// If ffmpeg is no longer running exit code != 0, it has finished
		if checkInspectResp.ExitCode != 0 {
			log.Printf("FFmpeg has finished writing video file for session %d", sessionId)
			break
		}

		log.Printf("FFmpeg still running, waiting... (attempt %d/10)", i+1)
		time.Sleep(1 * time.Second)
	}

	// Final check - if ffmpeg is still running after 10 seconds, force kill it
	forceKillCmd := []string{"pkill", "-SIGKILL", "ffmpeg"}
	forceKillExecResp, err := cli.ContainerExecCreate(context.Background(), containerID, container.ExecOptions{
		User:         "0:0",
		Cmd:          forceKillCmd,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          false,
		Detach:       false,
	})
	if err == nil {
		cli.ContainerExecStart(ctx, forceKillExecResp.ID, container.ExecStartOptions{})
		log.Printf("Force killed ffmpeg process for session %d", sessionId)
	}

	// Verify the video file was created and has content
	verifyCmd := []string{"bash", "-c", fmt.Sprintf("test -f /recordings/session_%d.mp4 && ls -la /recordings/session_%d.mp4", sessionId, sessionId)}
	verifyExecResp, err := cli.ContainerExecCreate(context.Background(), containerID, container.ExecOptions{
		User:         "0:0",
		Cmd:          verifyCmd,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          false,
		Detach:       false,
	})
	if err == nil {
		cli.ContainerExecStart(ctx, verifyExecResp.ID, container.ExecStartOptions{})
		verifyInspectResp, _ := cli.ContainerExecInspect(ctx, verifyExecResp.ID)
		if verifyInspectResp.ExitCode == 0 {
			log.Printf("Video file successfully created for session %d", sessionId)
		} else {
			log.Printf("Warning: Video file may not have been created properly for session %d", sessionId)
		}
	}

	log.Printf("Recording cleanup completed for session %d", sessionId)
	return nil
}

func startRecording(cli *client.Client, containerID string, sessionId int64) (string, error) {
	execErrCh := make(chan error, 1)
	ffmpegCmd := []string{
		"ffmpeg",
		"-f", "x11grab",
		"-video_size", "1360x1020",
		"-i", ":99", // The display source
		"-y",
		"/recordings/session_" + fmt.Sprintf("%d", sessionId) + ".mp4",
	}

	execConfig := container.ExecOptions{
		User:         "0:0",
		Cmd:          ffmpegCmd,
		Env:          []string{"DISPLAY=:99"},
		AttachStdout: true,
		AttachStderr: true,
		Tty:          false,
		Detach:       true,
	}

	execResp, err := cli.ContainerExecCreate(context.Background(), containerID, execConfig)
	if err != nil {
		log.Printf("Failed to create exec: %v", err)
		return "", err
	}
	log.Printf("Created exec: %+v", execResp)

	go func() {
		err := cli.ContainerExecStart(context.Background(), execResp.ID, container.ExecStartOptions{})
		execErrCh <- err
	}()

	select {
	case err := <-execErrCh:
		if err != nil {
			log.Printf("Error: Docker API call failed to start recording: %v", err)
			return "", fmt.Errorf("docker exec start failure: %w", err)
		}
		log.Printf("Success: Docker Exec process created and started for session %d", sessionId)
		return execResp.ID, nil
	}
}
