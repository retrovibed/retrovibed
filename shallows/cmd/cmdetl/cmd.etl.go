package cmdetl

type Commands struct {
	Completion cmdCompletion `cmd:"" help:"pipe jsonl prompts through a llama.cpp completion endpoint and emit results as jsonl"`
}
