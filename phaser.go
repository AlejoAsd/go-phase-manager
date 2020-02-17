package phaser

type PhaseManager interface {
	AddPhase(phaseName string, phase Phase)
	AddPreHookToPhase(phaseName string, phase Phase, hook PhaseHook)
	AddPostHookToPhase(phaseName string, phase Phase, hook PhaseHook)
}
