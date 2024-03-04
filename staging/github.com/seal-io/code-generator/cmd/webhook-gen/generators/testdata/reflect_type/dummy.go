package reflect_type

// DummyReconciler reconciles a Dummy object.
//
// nolint:lll
// +k8s:webhook-gen:validating:group="samplecontroller.seal.io",version="v1",resource="dummies",scope="*",operations=["CREATE","UPDATE","DELETE"],failurePolicy="Fail",sideEffects="None",matchPolicy="Equivalent",timeoutSeconds=10,namespaceSelector={"matchExpressions":[{"key":"walrus.seal.io/operate","operator":"Exists"}]},objectSelector={"matchLabels":{"label":"value"}},matchConditions=[{"name":"test","expression":"self.b > 0"}]
// +k8s:webhook-gen:mutating:group="samplecontroller.seal.io",version="v1",resource="dummies",scope="*",failurePolicy="Fail",sideEffects="None",matchPolicy="Equivalent",reinvocationPolicy="Never",timeoutSeconds=10,namespaceSelector={"matchExpressions":[{"key":"walrus.seal.io/operate","operator":"Exists"}]},objectSelector={"matchLabels":{"label":"value"}}
type DummyReconciler struct {
}
