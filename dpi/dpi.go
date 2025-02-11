// Package godpi provides the main API interface for utilizing the go-dpi library.
package dpi

import (
	"example.com/dpi/modules/classifiers"
	"github.com/google/gopacket"
	"github.com/mushorg/go-dpi/modules/classifiers"
	"github.com/mushorg/go-dpi/modules/ml"
	"github.com/mushorg/go-dpi/modules/wrappers"
	"github.com/mushorg/go-dpi/types"
	"time"
)

var activatedModules []types.Module
var moduleList = []types.Module{
	classifiers.NewClassifierModule(),
	wrappers.NewWrapperModule(),
	ml.NewLinearSVCModule(),
}
var cacheExpiration = 5 * time.Minute

func Initialize(opts ...Options) (errs []error) {

	for _, opt := range opts {
		if opt == nil {
			continue
		}

		for _, m := range moduleList {
			opt.Apply(m)
		}
	}
	types.InitCache(cacheExpiration)
	for _, module := range moduleList {
		activated := false
		for _, activeModule := range activatedModules {
			if activeModule == module {
				activated = true
				break
			}
		}
		if !activated {
			err := module.Initialize()
			if err == nil {
				activatedModules = append(activatedModules, module)
			} else {
				errs = append(errs, err)
			}
		}
	}
	return
}

func Destroy() (errs []error) {
	types.DestroyCache()
	newActivatedModules := make([]types.Module, 0)
	for _, module := range activatedModules {
		err := module.Destroy()
		if err != nil {
			newActivatedModules = append(newActivatedModules, module)
			errs = append(errs, err)
		}
	}
	activatedModules = newActivatedModules
	return
}

func SetModules(modules []types.Module) {
	moduleList = make([]types.Module, len(modules))
	copy(moduleList, modules)
}

func SetCacheExpiration(expiration time.Duration) {
	cacheExpiration = expiration
}

func GetPacketFlow(packet gopacket.Packet) (*types.Flow, bool) {
	return types.GetFlowForPacket(packet)
}

func ClassifyFlow(flow *types.Flow) (result types.ClassificationResult) {
	for _, module := range activatedModules {
		resultTmp := module.ClassifyFlow(flow)
		if resultTmp.Protocol != types.Unknown {
			result = resultTmp
			return
		}
	}
	return
}

func ClassifyFlowAllModules(flow *types.Flow) (results []types.ClassificationResult) {
	for _, module := range activatedModules {
		resultsTmp := module.ClassifyFlowAll(flow)
		for _, resultTmp := range resultsTmp {
			if resultTmp.Protocol != types.Unknown {
				results = append(results, resultTmp)
			}
		}
	}
	return
}
