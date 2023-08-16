// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// Code generated by "seal". DO NOT EDIT.

package model

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/seal-io/walrus/pkg/dao/model/costreport"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/object"
)

// CostReportCreateInput holds the creation input of the CostReport entity,
// please tags with `path:",inline" json:",inline"` if embedding.
type CostReportCreateInput struct {
	inputConfig `path:"-" query:"-" json:"-"`

	// String generated from resource properties, used to identify this cost.
	Fingerprint string `path:"-" query:"-" json:"fingerprint"`
	// Resource name for current cost, could be __unmounted__.
	Name string `path:"-" query:"-" json:"name"`
	// Usage minutes from start time to end time.
	Minutes float64 `path:"-" query:"-" json:"minutes"`
	// Usage end time for current cost.
	EndTime time.Time `path:"-" query:"-" json:"endTime"`
	// Usage start time for current cost.
	StartTime time.Time `path:"-" query:"-" json:"startTime"`
	// Cluster name for current cost.
	ClusterName string `path:"-" query:"-" json:"clusterName,omitempty"`
	// Namespace for current cost.
	Namespace string `path:"-" query:"-" json:"namespace,omitempty"`
	// Node for current cost.
	Node string `path:"-" query:"-" json:"node,omitempty"`
	// Controller name for the cost linked resource.
	Controller string `path:"-" query:"-" json:"controller,omitempty"`
	// Controller kind for the cost linked resource, deployment, statefulSet etc.
	ControllerKind string `path:"-" query:"-" json:"controllerKind,omitempty"`
	// Pod name for current cost.
	Pod string `path:"-" query:"-" json:"pod,omitempty"`
	// Container name for current cost.
	Container string `path:"-" query:"-" json:"container,omitempty"`
	// PV list for current cost linked.
	Pvs map[string]types.PVCost `path:"-" query:"-" json:"pvs,omitempty"`
	// Labels for the cost linked resource.
	Labels map[string]string `path:"-" query:"-" json:"labels,omitempty"`
	// Cost number.
	TotalCost float64 `path:"-" query:"-" json:"totalCost,omitempty"`
	// Cost currency.
	Currency int `path:"-" query:"-" json:"currency,omitempty"`
	// Cpu cost for current cost.
	CPUCost float64 `path:"-" query:"-" json:"cpuCost,omitempty"`
	// Cpu core requested.
	CPUCoreRequest float64 `path:"-" query:"-" json:"cpuCoreRequest,omitempty"`
	// GPU cost for current cost.
	GPUCost float64 `path:"-" query:"-" json:"gpuCost,omitempty"`
	// GPU core count.
	GPUCount float64 `path:"-" query:"-" json:"gpuCount,omitempty"`
	// Ram cost for current cost.
	RAMCost float64 `path:"-" query:"-" json:"ramCost,omitempty"`
	// Ram requested in byte.
	RAMByteRequest float64 `path:"-" query:"-" json:"ramByteRequest,omitempty"`
	// PV cost for current cost linked.
	PVCost float64 `path:"-" query:"-" json:"pvCost,omitempty"`
	// PV bytes for current cost linked.
	PVBytes float64 `path:"-" query:"-" json:"pvBytes,omitempty"`
	// LoadBalancer cost for current cost linked.
	LoadBalancerCost float64 `path:"-" query:"-" json:"loadBalancerCost,omitempty"`
	// CPU core average usage.
	CPUCoreUsageAverage float64 `path:"-" query:"-" json:"cpuCoreUsageAverage,omitempty"`
	// CPU core max usage.
	CPUCoreUsageMax float64 `path:"-" query:"-" json:"cpuCoreUsageMax,omitempty"`
	// Ram average usage in byte.
	RAMByteUsageAverage float64 `path:"-" query:"-" json:"ramByteUsageAverage,omitempty"`
	// Ram max usage in byte.
	RAMByteUsageMax float64 `path:"-" query:"-" json:"ramByteUsageMax,omitempty"`
}

// Model returns the CostReport entity for creating,
// after validating.
func (crci *CostReportCreateInput) Model() *CostReport {
	if crci == nil {
		return nil
	}

	_cr := &CostReport{
		Fingerprint:         crci.Fingerprint,
		Name:                crci.Name,
		Minutes:             crci.Minutes,
		EndTime:             crci.EndTime,
		StartTime:           crci.StartTime,
		ClusterName:         crci.ClusterName,
		Namespace:           crci.Namespace,
		Node:                crci.Node,
		Controller:          crci.Controller,
		ControllerKind:      crci.ControllerKind,
		Pod:                 crci.Pod,
		Container:           crci.Container,
		Pvs:                 crci.Pvs,
		Labels:              crci.Labels,
		TotalCost:           crci.TotalCost,
		Currency:            crci.Currency,
		CPUCost:             crci.CPUCost,
		CPUCoreRequest:      crci.CPUCoreRequest,
		GPUCost:             crci.GPUCost,
		GPUCount:            crci.GPUCount,
		RAMCost:             crci.RAMCost,
		RAMByteRequest:      crci.RAMByteRequest,
		PVCost:              crci.PVCost,
		PVBytes:             crci.PVBytes,
		LoadBalancerCost:    crci.LoadBalancerCost,
		CPUCoreUsageAverage: crci.CPUCoreUsageAverage,
		CPUCoreUsageMax:     crci.CPUCoreUsageMax,
		RAMByteUsageAverage: crci.RAMByteUsageAverage,
		RAMByteUsageMax:     crci.RAMByteUsageMax,
	}

	return _cr
}

// Validate checks the CostReportCreateInput entity.
func (crci *CostReportCreateInput) Validate() error {
	if crci == nil {
		return errors.New("nil receiver")
	}

	return crci.ValidateWith(crci.inputConfig.Context, crci.inputConfig.Client, nil)
}

// ValidateWith checks the CostReportCreateInput entity with the given context and client set.
func (crci *CostReportCreateInput) ValidateWith(ctx context.Context, cs ClientSet, cache map[string]any) error {
	if crci == nil {
		return errors.New("nil receiver")
	}

	if cache == nil {
		cache = map[string]any{}
	}

	return nil
}

// CostReportCreateInputs holds the creation input item of the CostReport entities.
type CostReportCreateInputsItem struct {
	// String generated from resource properties, used to identify this cost.
	Fingerprint string `path:"-" query:"-" json:"fingerprint"`
	// Resource name for current cost, could be __unmounted__.
	Name string `path:"-" query:"-" json:"name"`
	// Usage minutes from start time to end time.
	Minutes float64 `path:"-" query:"-" json:"minutes"`
	// Usage end time for current cost.
	EndTime time.Time `path:"-" query:"-" json:"endTime"`
	// Usage start time for current cost.
	StartTime time.Time `path:"-" query:"-" json:"startTime"`
	// Cluster name for current cost.
	ClusterName string `path:"-" query:"-" json:"clusterName,omitempty"`
	// Namespace for current cost.
	Namespace string `path:"-" query:"-" json:"namespace,omitempty"`
	// Node for current cost.
	Node string `path:"-" query:"-" json:"node,omitempty"`
	// Controller name for the cost linked resource.
	Controller string `path:"-" query:"-" json:"controller,omitempty"`
	// Controller kind for the cost linked resource, deployment, statefulSet etc.
	ControllerKind string `path:"-" query:"-" json:"controllerKind,omitempty"`
	// Pod name for current cost.
	Pod string `path:"-" query:"-" json:"pod,omitempty"`
	// Container name for current cost.
	Container string `path:"-" query:"-" json:"container,omitempty"`
	// PV list for current cost linked.
	Pvs map[string]types.PVCost `path:"-" query:"-" json:"pvs,omitempty"`
	// Labels for the cost linked resource.
	Labels map[string]string `path:"-" query:"-" json:"labels,omitempty"`
	// Cost number.
	TotalCost float64 `path:"-" query:"-" json:"totalCost,omitempty"`
	// Cost currency.
	Currency int `path:"-" query:"-" json:"currency,omitempty"`
	// Cpu cost for current cost.
	CPUCost float64 `path:"-" query:"-" json:"cpuCost,omitempty"`
	// Cpu core requested.
	CPUCoreRequest float64 `path:"-" query:"-" json:"cpuCoreRequest,omitempty"`
	// GPU cost for current cost.
	GPUCost float64 `path:"-" query:"-" json:"gpuCost,omitempty"`
	// GPU core count.
	GPUCount float64 `path:"-" query:"-" json:"gpuCount,omitempty"`
	// Ram cost for current cost.
	RAMCost float64 `path:"-" query:"-" json:"ramCost,omitempty"`
	// Ram requested in byte.
	RAMByteRequest float64 `path:"-" query:"-" json:"ramByteRequest,omitempty"`
	// PV cost for current cost linked.
	PVCost float64 `path:"-" query:"-" json:"pvCost,omitempty"`
	// PV bytes for current cost linked.
	PVBytes float64 `path:"-" query:"-" json:"pvBytes,omitempty"`
	// LoadBalancer cost for current cost linked.
	LoadBalancerCost float64 `path:"-" query:"-" json:"loadBalancerCost,omitempty"`
	// CPU core average usage.
	CPUCoreUsageAverage float64 `path:"-" query:"-" json:"cpuCoreUsageAverage,omitempty"`
	// CPU core max usage.
	CPUCoreUsageMax float64 `path:"-" query:"-" json:"cpuCoreUsageMax,omitempty"`
	// Ram average usage in byte.
	RAMByteUsageAverage float64 `path:"-" query:"-" json:"ramByteUsageAverage,omitempty"`
	// Ram max usage in byte.
	RAMByteUsageMax float64 `path:"-" query:"-" json:"ramByteUsageMax,omitempty"`
}

// ValidateWith checks the CostReportCreateInputsItem entity with the given context and client set.
func (crci *CostReportCreateInputsItem) ValidateWith(ctx context.Context, cs ClientSet, cache map[string]any) error {
	if crci == nil {
		return errors.New("nil receiver")
	}

	if cache == nil {
		cache = map[string]any{}
	}

	return nil
}

// CostReportCreateInputs holds the creation input of the CostReport entities,
// please tags with `path:",inline" json:",inline"` if embedding.
type CostReportCreateInputs struct {
	inputConfig `path:"-" query:"-" json:"-"`

	// Items holds the entities to create, which MUST not be empty.
	Items []*CostReportCreateInputsItem `path:"-" query:"-" json:"items"`
}

// Model returns the CostReport entities for creating,
// after validating.
func (crci *CostReportCreateInputs) Model() []*CostReport {
	if crci == nil || len(crci.Items) == 0 {
		return nil
	}

	_crs := make([]*CostReport, len(crci.Items))

	for i := range crci.Items {
		_cr := &CostReport{
			Fingerprint:         crci.Items[i].Fingerprint,
			Name:                crci.Items[i].Name,
			Minutes:             crci.Items[i].Minutes,
			EndTime:             crci.Items[i].EndTime,
			StartTime:           crci.Items[i].StartTime,
			ClusterName:         crci.Items[i].ClusterName,
			Namespace:           crci.Items[i].Namespace,
			Node:                crci.Items[i].Node,
			Controller:          crci.Items[i].Controller,
			ControllerKind:      crci.Items[i].ControllerKind,
			Pod:                 crci.Items[i].Pod,
			Container:           crci.Items[i].Container,
			Pvs:                 crci.Items[i].Pvs,
			Labels:              crci.Items[i].Labels,
			TotalCost:           crci.Items[i].TotalCost,
			Currency:            crci.Items[i].Currency,
			CPUCost:             crci.Items[i].CPUCost,
			CPUCoreRequest:      crci.Items[i].CPUCoreRequest,
			GPUCost:             crci.Items[i].GPUCost,
			GPUCount:            crci.Items[i].GPUCount,
			RAMCost:             crci.Items[i].RAMCost,
			RAMByteRequest:      crci.Items[i].RAMByteRequest,
			PVCost:              crci.Items[i].PVCost,
			PVBytes:             crci.Items[i].PVBytes,
			LoadBalancerCost:    crci.Items[i].LoadBalancerCost,
			CPUCoreUsageAverage: crci.Items[i].CPUCoreUsageAverage,
			CPUCoreUsageMax:     crci.Items[i].CPUCoreUsageMax,
			RAMByteUsageAverage: crci.Items[i].RAMByteUsageAverage,
			RAMByteUsageMax:     crci.Items[i].RAMByteUsageMax,
		}

		_crs[i] = _cr
	}

	return _crs
}

// Validate checks the CostReportCreateInputs entity .
func (crci *CostReportCreateInputs) Validate() error {
	if crci == nil {
		return errors.New("nil receiver")
	}

	return crci.ValidateWith(crci.inputConfig.Context, crci.inputConfig.Client, nil)
}

// ValidateWith checks the CostReportCreateInputs entity with the given context and client set.
func (crci *CostReportCreateInputs) ValidateWith(ctx context.Context, cs ClientSet, cache map[string]any) error {
	if crci == nil {
		return errors.New("nil receiver")
	}

	if len(crci.Items) == 0 {
		return errors.New("empty items")
	}

	if cache == nil {
		cache = map[string]any{}
	}

	for i := range crci.Items {
		if crci.Items[i] == nil {
			continue
		}

		if err := crci.Items[i].ValidateWith(ctx, cs, cache); err != nil {
			return err
		}
	}

	return nil
}

// CostReportDeleteInput holds the deletion input of the CostReport entity,
// please tags with `path:",inline"` if embedding.
type CostReportDeleteInput struct {
	CostReportQueryInput `path:",inline"`
}

// CostReportDeleteInputs holds the deletion input item of the CostReport entities.
type CostReportDeleteInputsItem struct {
	// ID of the CostReport entity.
	ID int `path:"-" query:"-" json:"id"`
}

// CostReportDeleteInputs holds the deletion input of the CostReport entities,
// please tags with `path:",inline" json:",inline"` if embedding.
type CostReportDeleteInputs struct {
	inputConfig `path:"-" query:"-" json:"-"`

	// Items holds the entities to create, which MUST not be empty.
	Items []*CostReportDeleteInputsItem `path:"-" query:"-" json:"items"`
}

// Model returns the CostReport entities for deleting,
// after validating.
func (crdi *CostReportDeleteInputs) Model() []*CostReport {
	if crdi == nil || len(crdi.Items) == 0 {
		return nil
	}

	_crs := make([]*CostReport, len(crdi.Items))
	for i := range crdi.Items {
		_crs[i] = &CostReport{
			ID: crdi.Items[i].ID,
		}
	}
	return _crs
}

// IDs returns the ID list of the CostReport entities for deleting,
// after validating.
func (crdi *CostReportDeleteInputs) IDs() []int {
	if crdi == nil || len(crdi.Items) == 0 {
		return nil
	}

	ids := make([]int, len(crdi.Items))
	for i := range crdi.Items {
		ids[i] = crdi.Items[i].ID
	}
	return ids
}

// Validate checks the CostReportDeleteInputs entity.
func (crdi *CostReportDeleteInputs) Validate() error {
	if crdi == nil {
		return errors.New("nil receiver")
	}

	return crdi.ValidateWith(crdi.inputConfig.Context, crdi.inputConfig.Client, nil)
}

// ValidateWith checks the CostReportDeleteInputs entity with the given context and client set.
func (crdi *CostReportDeleteInputs) ValidateWith(ctx context.Context, cs ClientSet, cache map[string]any) error {
	if crdi == nil {
		return errors.New("nil receiver")
	}

	if len(crdi.Items) == 0 {
		return errors.New("empty items")
	}

	if cache == nil {
		cache = map[string]any{}
	}

	q := cs.CostReports().Query()

	ids := make([]int, 0, len(crdi.Items))

	for i := range crdi.Items {
		if crdi.Items[i] == nil {
			return errors.New("nil item")
		}

		if crdi.Items[i].ID != 0 {
			ids = append(ids, crdi.Items[i].ID)
		} else {
			return errors.New("found item hasn't identify")
		}
	}

	if len(ids) != cap(ids) {
		return errors.New("found unrecognized item")
	}

	idsCnt, err := q.Where(costreport.IDIn(ids...)).
		Count(ctx)
	if err != nil {
		return err
	}

	if idsCnt != cap(ids) {
		return errors.New("found unrecognized item")
	}

	return nil
}

// CostReportQueryInput holds the query input of the CostReport entity,
// please tags with `path:",inline"` if embedding.
type CostReportQueryInput struct {
	inputConfig `path:"-" query:"-" json:"-"`

	// Refer holds the route path reference of the CostReport entity.
	Refer *object.Refer `path:"costreport,default=" query:"-" json:"-"`
	// ID of the CostReport entity.
	ID int `path:"-" query:"-" json:"id"`
}

// Model returns the CostReport entity for querying,
// after validating.
func (crqi *CostReportQueryInput) Model() *CostReport {
	if crqi == nil {
		return nil
	}

	return &CostReport{
		ID: crqi.ID,
	}
}

// Validate checks the CostReportQueryInput entity.
func (crqi *CostReportQueryInput) Validate() error {
	if crqi == nil {
		return errors.New("nil receiver")
	}

	return crqi.ValidateWith(crqi.inputConfig.Context, crqi.inputConfig.Client, nil)
}

// ValidateWith checks the CostReportQueryInput entity with the given context and client set.
func (crqi *CostReportQueryInput) ValidateWith(ctx context.Context, cs ClientSet, cache map[string]any) error {
	if crqi == nil {
		return errors.New("nil receiver")
	}

	if crqi.Refer != nil && *crqi.Refer == "" {
		return fmt.Errorf("model: %s : %w", costreport.Label, ErrBlankResourceRefer)
	}

	if cache == nil {
		cache = map[string]any{}
	}

	q := cs.CostReports().Query()

	if crqi.Refer != nil {
		if crqi.Refer.IsNumeric() {
			q.Where(
				costreport.ID(crqi.Refer.Int()))
		} else {
			return errors.New("invalid identify refer of costreport")
		}
	} else if crqi.ID != 0 {
		q.Where(
			costreport.ID(crqi.ID))
	} else {
		return errors.New("invalid identify of costreport")
	}

	q.Select(
		costreport.FieldID,
	)

	var e *CostReport
	{
		// Get cache from previous validation.
		queryStmt, queryArgs := q.sqlQuery(setContextOp(ctx, q.ctx, "cache")).Query()
		ck := fmt.Sprintf("stmt=%v, args=%v", queryStmt, queryArgs)
		if cv, existed := cache[ck]; !existed {
			var err error
			e, err = q.Only(ctx)
			if err != nil {
				return err
			}

			// Set cache for other validation.
			cache[ck] = e
		} else {
			e = cv.(*CostReport)
		}
	}

	crqi.ID = e.ID
	return nil
}

// CostReportQueryInputs holds the query input of the CostReport entities,
// please tags with `path:",inline" query:",inline"` if embedding.
type CostReportQueryInputs struct {
	inputConfig `path:"-" query:"-" json:"-"`
}

// Validate checks the CostReportQueryInputs entity.
func (crqi *CostReportQueryInputs) Validate() error {
	if crqi == nil {
		return errors.New("nil receiver")
	}

	return crqi.ValidateWith(crqi.inputConfig.Context, crqi.inputConfig.Client, nil)
}

// ValidateWith checks the CostReportQueryInputs entity with the given context and client set.
func (crqi *CostReportQueryInputs) ValidateWith(ctx context.Context, cs ClientSet, cache map[string]any) error {
	if crqi == nil {
		return errors.New("nil receiver")
	}

	if cache == nil {
		cache = map[string]any{}
	}

	return nil
}

// CostReportUpdateInput holds the modification input of the CostReport entity,
// please tags with `path:",inline" json:",inline"` if embedding.
type CostReportUpdateInput struct {
	CostReportQueryInput `path:",inline" query:"-" json:"-"`

	// Cost number.
	TotalCost float64 `path:"-" query:"-" json:"totalCost,omitempty"`
	// Cost currency.
	Currency int `path:"-" query:"-" json:"currency,omitempty"`
	// Cpu cost for current cost.
	CPUCost float64 `path:"-" query:"-" json:"cpuCost,omitempty"`
	// GPU cost for current cost.
	GPUCost float64 `path:"-" query:"-" json:"gpuCost,omitempty"`
	// Ram cost for current cost.
	RAMCost float64 `path:"-" query:"-" json:"ramCost,omitempty"`
	// PV cost for current cost linked.
	PVCost float64 `path:"-" query:"-" json:"pvCost,omitempty"`
	// PV bytes for current cost linked.
	PVBytes float64 `path:"-" query:"-" json:"pvBytes,omitempty"`
	// LoadBalancer cost for current cost linked.
	LoadBalancerCost float64 `path:"-" query:"-" json:"loadBalancerCost,omitempty"`
	// CPU core average usage.
	CPUCoreUsageAverage float64 `path:"-" query:"-" json:"cpuCoreUsageAverage,omitempty"`
	// CPU core max usage.
	CPUCoreUsageMax float64 `path:"-" query:"-" json:"cpuCoreUsageMax,omitempty"`
	// Ram average usage in byte.
	RAMByteUsageAverage float64 `path:"-" query:"-" json:"rambyteUsageAverage,omitempty"`
	// Ram max usage in byte.
	RAMByteUsageMax float64 `path:"-" query:"-" json:"rambyteUsageMax,omitempty"`
}

// Model returns the CostReport entity for modifying,
// after validating.
func (crui *CostReportUpdateInput) Model() *CostReport {
	if crui == nil {
		return nil
	}

	_cr := &CostReport{
		ID:                  crui.ID,
		TotalCost:           crui.TotalCost,
		Currency:            crui.Currency,
		CPUCost:             crui.CPUCost,
		GPUCost:             crui.GPUCost,
		RAMCost:             crui.RAMCost,
		PVCost:              crui.PVCost,
		PVBytes:             crui.PVBytes,
		LoadBalancerCost:    crui.LoadBalancerCost,
		CPUCoreUsageAverage: crui.CPUCoreUsageAverage,
		CPUCoreUsageMax:     crui.CPUCoreUsageMax,
		RAMByteUsageAverage: crui.RAMByteUsageAverage,
		RAMByteUsageMax:     crui.RAMByteUsageMax,
	}

	return _cr
}

// Validate checks the CostReportUpdateInput entity.
func (crui *CostReportUpdateInput) Validate() error {
	if crui == nil {
		return errors.New("nil receiver")
	}

	return crui.ValidateWith(crui.inputConfig.Context, crui.inputConfig.Client, nil)
}

// ValidateWith checks the CostReportUpdateInput entity with the given context and client set.
func (crui *CostReportUpdateInput) ValidateWith(ctx context.Context, cs ClientSet, cache map[string]any) error {
	if cache == nil {
		cache = map[string]any{}
	}

	if err := crui.CostReportQueryInput.ValidateWith(ctx, cs, cache); err != nil {
		return err
	}

	return nil
}

// CostReportUpdateInputs holds the modification input item of the CostReport entities.
type CostReportUpdateInputsItem struct {
	// ID of the CostReport entity.
	ID int `path:"-" query:"-" json:"id"`

	// Cost number.
	TotalCost float64 `path:"-" query:"-" json:"totalCost"`
	// Cost currency.
	Currency int `path:"-" query:"-" json:"currency,omitempty"`
	// Cpu cost for current cost.
	CPUCost float64 `path:"-" query:"-" json:"cpuCost"`
	// GPU cost for current cost.
	GPUCost float64 `path:"-" query:"-" json:"gpuCost"`
	// Ram cost for current cost.
	RAMCost float64 `path:"-" query:"-" json:"ramCost"`
	// PV cost for current cost linked.
	PVCost float64 `path:"-" query:"-" json:"pvCost"`
	// PV bytes for current cost linked.
	PVBytes float64 `path:"-" query:"-" json:"pvBytes"`
	// LoadBalancer cost for current cost linked.
	LoadBalancerCost float64 `path:"-" query:"-" json:"loadBalancerCost"`
	// CPU core average usage.
	CPUCoreUsageAverage float64 `path:"-" query:"-" json:"cpuCoreUsageAverage"`
	// CPU core max usage.
	CPUCoreUsageMax float64 `path:"-" query:"-" json:"cpuCoreUsageMax"`
	// Ram average usage in byte.
	RAMByteUsageAverage float64 `path:"-" query:"-" json:"ramByteUsageAverage"`
	// Ram max usage in byte.
	RAMByteUsageMax float64 `path:"-" query:"-" json:"ramByteUsageMax"`
}

// ValidateWith checks the CostReportUpdateInputsItem entity with the given context and client set.
func (crui *CostReportUpdateInputsItem) ValidateWith(ctx context.Context, cs ClientSet, cache map[string]any) error {
	if crui == nil {
		return errors.New("nil receiver")
	}

	if cache == nil {
		cache = map[string]any{}
	}

	return nil
}

// CostReportUpdateInputs holds the modification input of the CostReport entities,
// please tags with `path:",inline" json:",inline"` if embedding.
type CostReportUpdateInputs struct {
	inputConfig `path:"-" query:"-" json:"-"`

	// Items holds the entities to create, which MUST not be empty.
	Items []*CostReportUpdateInputsItem `path:"-" query:"-" json:"items"`
}

// Model returns the CostReport entities for modifying,
// after validating.
func (crui *CostReportUpdateInputs) Model() []*CostReport {
	if crui == nil || len(crui.Items) == 0 {
		return nil
	}

	_crs := make([]*CostReport, len(crui.Items))

	for i := range crui.Items {
		_cr := &CostReport{
			ID:                  crui.Items[i].ID,
			TotalCost:           crui.Items[i].TotalCost,
			Currency:            crui.Items[i].Currency,
			CPUCost:             crui.Items[i].CPUCost,
			GPUCost:             crui.Items[i].GPUCost,
			RAMCost:             crui.Items[i].RAMCost,
			PVCost:              crui.Items[i].PVCost,
			PVBytes:             crui.Items[i].PVBytes,
			LoadBalancerCost:    crui.Items[i].LoadBalancerCost,
			CPUCoreUsageAverage: crui.Items[i].CPUCoreUsageAverage,
			CPUCoreUsageMax:     crui.Items[i].CPUCoreUsageMax,
			RAMByteUsageAverage: crui.Items[i].RAMByteUsageAverage,
			RAMByteUsageMax:     crui.Items[i].RAMByteUsageMax,
		}

		_crs[i] = _cr
	}

	return _crs
}

// IDs returns the ID list of the CostReport entities for modifying,
// after validating.
func (crui *CostReportUpdateInputs) IDs() []int {
	if crui == nil || len(crui.Items) == 0 {
		return nil
	}

	ids := make([]int, len(crui.Items))
	for i := range crui.Items {
		ids[i] = crui.Items[i].ID
	}
	return ids
}

// Validate checks the CostReportUpdateInputs entity.
func (crui *CostReportUpdateInputs) Validate() error {
	if crui == nil {
		return errors.New("nil receiver")
	}

	return crui.ValidateWith(crui.inputConfig.Context, crui.inputConfig.Client, nil)
}

// ValidateWith checks the CostReportUpdateInputs entity with the given context and client set.
func (crui *CostReportUpdateInputs) ValidateWith(ctx context.Context, cs ClientSet, cache map[string]any) error {
	if crui == nil {
		return errors.New("nil receiver")
	}

	if len(crui.Items) == 0 {
		return errors.New("empty items")
	}

	if cache == nil {
		cache = map[string]any{}
	}

	q := cs.CostReports().Query()

	ids := make([]int, 0, len(crui.Items))

	for i := range crui.Items {
		if crui.Items[i] == nil {
			return errors.New("nil item")
		}

		if crui.Items[i].ID != 0 {
			ids = append(ids, crui.Items[i].ID)
		} else {
			return errors.New("found item hasn't identify")
		}
	}

	if len(ids) != cap(ids) {
		return errors.New("found unrecognized item")
	}

	idsCnt, err := q.Where(costreport.IDIn(ids...)).
		Count(ctx)
	if err != nil {
		return err
	}

	if idsCnt != cap(ids) {
		return errors.New("found unrecognized item")
	}

	for i := range crui.Items {
		if crui.Items[i] == nil {
			continue
		}

		if err := crui.Items[i].ValidateWith(ctx, cs, cache); err != nil {
			return err
		}
	}

	return nil
}

// CostReportOutput holds the output of the CostReport entity.
type CostReportOutput struct {
	ID                  int                     `json:"id,omitempty"`
	StartTime           time.Time               `json:"startTime,omitempty"`
	EndTime             time.Time               `json:"endTime,omitempty"`
	Minutes             float64                 `json:"minutes,omitempty"`
	Name                string                  `json:"name,omitempty"`
	Fingerprint         string                  `json:"fingerprint,omitempty"`
	ClusterName         string                  `json:"clusterName,omitempty"`
	Namespace           string                  `json:"namespace,omitempty"`
	Node                string                  `json:"node,omitempty"`
	Controller          string                  `json:"controller,omitempty"`
	ControllerKind      string                  `json:"controllerKind,omitempty"`
	Pod                 string                  `json:"pod,omitempty"`
	Container           string                  `json:"container,omitempty"`
	Pvs                 map[string]types.PVCost `json:"pvs,omitempty"`
	Labels              map[string]string       `json:"labels,omitempty"`
	TotalCost           float64                 `json:"totalCost,omitempty"`
	Currency            int                     `json:"currency,omitempty"`
	CPUCost             float64                 `json:"cpuCost,omitempty"`
	CPUCoreRequest      float64                 `json:"cpuCoreRequest,omitempty"`
	GPUCost             float64                 `json:"gpuCost,omitempty"`
	GPUCount            float64                 `json:"gpuCount,omitempty"`
	RAMCost             float64                 `json:"ramCost,omitempty"`
	RAMByteRequest      float64                 `json:"ramByteRequest,omitempty"`
	PVCost              float64                 `json:"pvCost,omitempty"`
	PVBytes             float64                 `json:"pvBytes,omitempty"`
	LoadBalancerCost    float64                 `json:"loadBalancerCost,omitempty"`
	CPUCoreUsageAverage float64                 `json:"cpuCoreUsageAverage,omitempty"`
	CPUCoreUsageMax     float64                 `json:"cpuCoreUsageMax,omitempty"`
	RAMByteUsageAverage float64                 `json:"ramByteUsageAverage,omitempty"`
	RAMByteUsageMax     float64                 `json:"ramByteUsageMax,omitempty"`
}

// View returns the output of CostReport entity.
func (_cr *CostReport) View() *CostReportOutput {
	return ExposeCostReport(_cr)
}

// View returns the output of CostReport entities.
func (_crs CostReports) View() []*CostReportOutput {
	return ExposeCostReports(_crs)
}

// ExposeCostReport converts the CostReport to CostReportOutput.
func ExposeCostReport(_cr *CostReport) *CostReportOutput {
	if _cr == nil {
		return nil
	}

	cro := &CostReportOutput{
		ID:                  _cr.ID,
		StartTime:           _cr.StartTime,
		EndTime:             _cr.EndTime,
		Minutes:             _cr.Minutes,
		Name:                _cr.Name,
		Fingerprint:         _cr.Fingerprint,
		ClusterName:         _cr.ClusterName,
		Namespace:           _cr.Namespace,
		Node:                _cr.Node,
		Controller:          _cr.Controller,
		ControllerKind:      _cr.ControllerKind,
		Pod:                 _cr.Pod,
		Container:           _cr.Container,
		Pvs:                 _cr.Pvs,
		Labels:              _cr.Labels,
		TotalCost:           _cr.TotalCost,
		Currency:            _cr.Currency,
		CPUCost:             _cr.CPUCost,
		CPUCoreRequest:      _cr.CPUCoreRequest,
		GPUCost:             _cr.GPUCost,
		GPUCount:            _cr.GPUCount,
		RAMCost:             _cr.RAMCost,
		RAMByteRequest:      _cr.RAMByteRequest,
		PVCost:              _cr.PVCost,
		PVBytes:             _cr.PVBytes,
		LoadBalancerCost:    _cr.LoadBalancerCost,
		CPUCoreUsageAverage: _cr.CPUCoreUsageAverage,
		CPUCoreUsageMax:     _cr.CPUCoreUsageMax,
		RAMByteUsageAverage: _cr.RAMByteUsageAverage,
		RAMByteUsageMax:     _cr.RAMByteUsageMax,
	}

	return cro
}

// ExposeCostReports converts the CostReport slice to CostReportOutput pointer slice.
func ExposeCostReports(_crs []*CostReport) []*CostReportOutput {
	if len(_crs) == 0 {
		return nil
	}

	cros := make([]*CostReportOutput, len(_crs))
	for i := range _crs {
		cros[i] = ExposeCostReport(_crs[i])
	}
	return cros
}
