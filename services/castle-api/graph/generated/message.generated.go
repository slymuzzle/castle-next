// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package generated

import (
	"context"
	"errors"
	"fmt"
	"io"
	"journeyhub/ent"
	"journeyhub/ent/schema/pulid"
	"journeyhub/graph/model"
	"strconv"

	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/v2/ast"
)

// region    ************************** generated!.gotpl **************************

type SubscriptionResolver interface {
	MessageCreated(ctx context.Context, roomID pulid.ID) (<-chan *ent.MessageEdge, error)
	MessageUpdated(ctx context.Context, roomID pulid.ID) (<-chan *ent.MessageEdge, error)
	MessageDeleted(ctx context.Context, roomID pulid.ID) (<-chan pulid.ID, error)
	RoomMemberCreated(ctx context.Context) (<-chan *ent.RoomMemberEdge, error)
	RoomMemberUpdated(ctx context.Context) (<-chan *ent.RoomMemberEdge, error)
	RoomMemberDeleted(ctx context.Context) (<-chan pulid.ID, error)
}

// endregion ************************** generated!.gotpl **************************

// region    ***************************** args.gotpl *****************************

func (ec *executionContext) field_Subscription_messageCreated_args(ctx context.Context, rawArgs map[string]interface{}) (map[string]interface{}, error) {
	var err error
	args := map[string]interface{}{}
	var arg0 pulid.ID
	if tmp, ok := rawArgs["roomID"]; ok {
		ctx := graphql.WithPathContext(ctx, graphql.NewPathWithField("roomID"))
		arg0, err = ec.unmarshalNID2journeyhubᚋentᚋschemaᚋpulidᚐID(ctx, tmp)
		if err != nil {
			return nil, err
		}
	}
	args["roomID"] = arg0
	return args, nil
}

func (ec *executionContext) field_Subscription_messageDeleted_args(ctx context.Context, rawArgs map[string]interface{}) (map[string]interface{}, error) {
	var err error
	args := map[string]interface{}{}
	var arg0 pulid.ID
	if tmp, ok := rawArgs["roomID"]; ok {
		ctx := graphql.WithPathContext(ctx, graphql.NewPathWithField("roomID"))
		arg0, err = ec.unmarshalNID2journeyhubᚋentᚋschemaᚋpulidᚐID(ctx, tmp)
		if err != nil {
			return nil, err
		}
	}
	args["roomID"] = arg0
	return args, nil
}

func (ec *executionContext) field_Subscription_messageUpdated_args(ctx context.Context, rawArgs map[string]interface{}) (map[string]interface{}, error) {
	var err error
	args := map[string]interface{}{}
	var arg0 pulid.ID
	if tmp, ok := rawArgs["roomID"]; ok {
		ctx := graphql.WithPathContext(ctx, graphql.NewPathWithField("roomID"))
		arg0, err = ec.unmarshalNID2journeyhubᚋentᚋschemaᚋpulidᚐID(ctx, tmp)
		if err != nil {
			return nil, err
		}
	}
	args["roomID"] = arg0
	return args, nil
}

// endregion ***************************** args.gotpl *****************************

// region    ************************** directives.gotpl **************************

// endregion ************************** directives.gotpl **************************

// region    **************************** field.gotpl *****************************

func (ec *executionContext) _Subscription_messageCreated(ctx context.Context, field graphql.CollectedField) (ret func(ctx context.Context) graphql.Marshaler) {
	fc, err := ec.fieldContext_Subscription_messageCreated(ctx, field)
	if err != nil {
		return nil
	}
	ctx = graphql.WithFieldContext(ctx, fc)
	defer func() {
		if r := recover(); r != nil {
			ec.Error(ctx, ec.Recover(ctx, r))
			ret = nil
		}
	}()
	resTmp, err := ec.ResolverMiddleware(ctx, func(rctx context.Context) (interface{}, error) {
		ctx = rctx // use context from middleware stack in children
		return ec.resolvers.Subscription().MessageCreated(rctx, fc.Args["roomID"].(pulid.ID))
	})
	if err != nil {
		ec.Error(ctx, err)
		return nil
	}
	if resTmp == nil {
		if !graphql.HasFieldError(ctx, fc) {
			ec.Errorf(ctx, "must not be null")
		}
		return nil
	}
	return func(ctx context.Context) graphql.Marshaler {
		select {
		case res, ok := <-resTmp.(<-chan *ent.MessageEdge):
			if !ok {
				return nil
			}
			return graphql.WriterFunc(func(w io.Writer) {
				w.Write([]byte{'{'})
				graphql.MarshalString(field.Alias).MarshalGQL(w)
				w.Write([]byte{':'})
				ec.marshalNMessageEdge2ᚖjourneyhubᚋentᚐMessageEdge(ctx, field.Selections, res).MarshalGQL(w)
				w.Write([]byte{'}'})
			})
		case <-ctx.Done():
			return nil
		}
	}
}

func (ec *executionContext) fieldContext_Subscription_messageCreated(ctx context.Context, field graphql.CollectedField) (fc *graphql.FieldContext, err error) {
	fc = &graphql.FieldContext{
		Object:     "Subscription",
		Field:      field,
		IsMethod:   true,
		IsResolver: true,
		Child: func(ctx context.Context, field graphql.CollectedField) (*graphql.FieldContext, error) {
			switch field.Name {
			case "node":
				return ec.fieldContext_MessageEdge_node(ctx, field)
			case "cursor":
				return ec.fieldContext_MessageEdge_cursor(ctx, field)
			}
			return nil, fmt.Errorf("no field named %q was found under type MessageEdge", field.Name)
		},
	}
	defer func() {
		if r := recover(); r != nil {
			err = ec.Recover(ctx, r)
			ec.Error(ctx, err)
		}
	}()
	ctx = graphql.WithFieldContext(ctx, fc)
	if fc.Args, err = ec.field_Subscription_messageCreated_args(ctx, field.ArgumentMap(ec.Variables)); err != nil {
		ec.Error(ctx, err)
		return fc, err
	}
	return fc, nil
}

func (ec *executionContext) _Subscription_messageUpdated(ctx context.Context, field graphql.CollectedField) (ret func(ctx context.Context) graphql.Marshaler) {
	fc, err := ec.fieldContext_Subscription_messageUpdated(ctx, field)
	if err != nil {
		return nil
	}
	ctx = graphql.WithFieldContext(ctx, fc)
	defer func() {
		if r := recover(); r != nil {
			ec.Error(ctx, ec.Recover(ctx, r))
			ret = nil
		}
	}()
	resTmp, err := ec.ResolverMiddleware(ctx, func(rctx context.Context) (interface{}, error) {
		ctx = rctx // use context from middleware stack in children
		return ec.resolvers.Subscription().MessageUpdated(rctx, fc.Args["roomID"].(pulid.ID))
	})
	if err != nil {
		ec.Error(ctx, err)
		return nil
	}
	if resTmp == nil {
		if !graphql.HasFieldError(ctx, fc) {
			ec.Errorf(ctx, "must not be null")
		}
		return nil
	}
	return func(ctx context.Context) graphql.Marshaler {
		select {
		case res, ok := <-resTmp.(<-chan *ent.MessageEdge):
			if !ok {
				return nil
			}
			return graphql.WriterFunc(func(w io.Writer) {
				w.Write([]byte{'{'})
				graphql.MarshalString(field.Alias).MarshalGQL(w)
				w.Write([]byte{':'})
				ec.marshalNMessageEdge2ᚖjourneyhubᚋentᚐMessageEdge(ctx, field.Selections, res).MarshalGQL(w)
				w.Write([]byte{'}'})
			})
		case <-ctx.Done():
			return nil
		}
	}
}

func (ec *executionContext) fieldContext_Subscription_messageUpdated(ctx context.Context, field graphql.CollectedField) (fc *graphql.FieldContext, err error) {
	fc = &graphql.FieldContext{
		Object:     "Subscription",
		Field:      field,
		IsMethod:   true,
		IsResolver: true,
		Child: func(ctx context.Context, field graphql.CollectedField) (*graphql.FieldContext, error) {
			switch field.Name {
			case "node":
				return ec.fieldContext_MessageEdge_node(ctx, field)
			case "cursor":
				return ec.fieldContext_MessageEdge_cursor(ctx, field)
			}
			return nil, fmt.Errorf("no field named %q was found under type MessageEdge", field.Name)
		},
	}
	defer func() {
		if r := recover(); r != nil {
			err = ec.Recover(ctx, r)
			ec.Error(ctx, err)
		}
	}()
	ctx = graphql.WithFieldContext(ctx, fc)
	if fc.Args, err = ec.field_Subscription_messageUpdated_args(ctx, field.ArgumentMap(ec.Variables)); err != nil {
		ec.Error(ctx, err)
		return fc, err
	}
	return fc, nil
}

func (ec *executionContext) _Subscription_messageDeleted(ctx context.Context, field graphql.CollectedField) (ret func(ctx context.Context) graphql.Marshaler) {
	fc, err := ec.fieldContext_Subscription_messageDeleted(ctx, field)
	if err != nil {
		return nil
	}
	ctx = graphql.WithFieldContext(ctx, fc)
	defer func() {
		if r := recover(); r != nil {
			ec.Error(ctx, ec.Recover(ctx, r))
			ret = nil
		}
	}()
	resTmp, err := ec.ResolverMiddleware(ctx, func(rctx context.Context) (interface{}, error) {
		ctx = rctx // use context from middleware stack in children
		return ec.resolvers.Subscription().MessageDeleted(rctx, fc.Args["roomID"].(pulid.ID))
	})
	if err != nil {
		ec.Error(ctx, err)
		return nil
	}
	if resTmp == nil {
		if !graphql.HasFieldError(ctx, fc) {
			ec.Errorf(ctx, "must not be null")
		}
		return nil
	}
	return func(ctx context.Context) graphql.Marshaler {
		select {
		case res, ok := <-resTmp.(<-chan pulid.ID):
			if !ok {
				return nil
			}
			return graphql.WriterFunc(func(w io.Writer) {
				w.Write([]byte{'{'})
				graphql.MarshalString(field.Alias).MarshalGQL(w)
				w.Write([]byte{':'})
				ec.marshalNID2journeyhubᚋentᚋschemaᚋpulidᚐID(ctx, field.Selections, res).MarshalGQL(w)
				w.Write([]byte{'}'})
			})
		case <-ctx.Done():
			return nil
		}
	}
}

func (ec *executionContext) fieldContext_Subscription_messageDeleted(ctx context.Context, field graphql.CollectedField) (fc *graphql.FieldContext, err error) {
	fc = &graphql.FieldContext{
		Object:     "Subscription",
		Field:      field,
		IsMethod:   true,
		IsResolver: true,
		Child: func(ctx context.Context, field graphql.CollectedField) (*graphql.FieldContext, error) {
			return nil, errors.New("field of type ID does not have child fields")
		},
	}
	defer func() {
		if r := recover(); r != nil {
			err = ec.Recover(ctx, r)
			ec.Error(ctx, err)
		}
	}()
	ctx = graphql.WithFieldContext(ctx, fc)
	if fc.Args, err = ec.field_Subscription_messageDeleted_args(ctx, field.ArgumentMap(ec.Variables)); err != nil {
		ec.Error(ctx, err)
		return fc, err
	}
	return fc, nil
}

func (ec *executionContext) _Subscription_roomMemberCreated(ctx context.Context, field graphql.CollectedField) (ret func(ctx context.Context) graphql.Marshaler) {
	fc, err := ec.fieldContext_Subscription_roomMemberCreated(ctx, field)
	if err != nil {
		return nil
	}
	ctx = graphql.WithFieldContext(ctx, fc)
	defer func() {
		if r := recover(); r != nil {
			ec.Error(ctx, ec.Recover(ctx, r))
			ret = nil
		}
	}()
	resTmp, err := ec.ResolverMiddleware(ctx, func(rctx context.Context) (interface{}, error) {
		ctx = rctx // use context from middleware stack in children
		return ec.resolvers.Subscription().RoomMemberCreated(rctx)
	})
	if err != nil {
		ec.Error(ctx, err)
		return nil
	}
	if resTmp == nil {
		if !graphql.HasFieldError(ctx, fc) {
			ec.Errorf(ctx, "must not be null")
		}
		return nil
	}
	return func(ctx context.Context) graphql.Marshaler {
		select {
		case res, ok := <-resTmp.(<-chan *ent.RoomMemberEdge):
			if !ok {
				return nil
			}
			return graphql.WriterFunc(func(w io.Writer) {
				w.Write([]byte{'{'})
				graphql.MarshalString(field.Alias).MarshalGQL(w)
				w.Write([]byte{':'})
				ec.marshalNRoomMemberEdge2ᚖjourneyhubᚋentᚐRoomMemberEdge(ctx, field.Selections, res).MarshalGQL(w)
				w.Write([]byte{'}'})
			})
		case <-ctx.Done():
			return nil
		}
	}
}

func (ec *executionContext) fieldContext_Subscription_roomMemberCreated(_ context.Context, field graphql.CollectedField) (fc *graphql.FieldContext, err error) {
	fc = &graphql.FieldContext{
		Object:     "Subscription",
		Field:      field,
		IsMethod:   true,
		IsResolver: true,
		Child: func(ctx context.Context, field graphql.CollectedField) (*graphql.FieldContext, error) {
			switch field.Name {
			case "node":
				return ec.fieldContext_RoomMemberEdge_node(ctx, field)
			case "cursor":
				return ec.fieldContext_RoomMemberEdge_cursor(ctx, field)
			}
			return nil, fmt.Errorf("no field named %q was found under type RoomMemberEdge", field.Name)
		},
	}
	return fc, nil
}

func (ec *executionContext) _Subscription_roomMemberUpdated(ctx context.Context, field graphql.CollectedField) (ret func(ctx context.Context) graphql.Marshaler) {
	fc, err := ec.fieldContext_Subscription_roomMemberUpdated(ctx, field)
	if err != nil {
		return nil
	}
	ctx = graphql.WithFieldContext(ctx, fc)
	defer func() {
		if r := recover(); r != nil {
			ec.Error(ctx, ec.Recover(ctx, r))
			ret = nil
		}
	}()
	resTmp, err := ec.ResolverMiddleware(ctx, func(rctx context.Context) (interface{}, error) {
		ctx = rctx // use context from middleware stack in children
		return ec.resolvers.Subscription().RoomMemberUpdated(rctx)
	})
	if err != nil {
		ec.Error(ctx, err)
		return nil
	}
	if resTmp == nil {
		if !graphql.HasFieldError(ctx, fc) {
			ec.Errorf(ctx, "must not be null")
		}
		return nil
	}
	return func(ctx context.Context) graphql.Marshaler {
		select {
		case res, ok := <-resTmp.(<-chan *ent.RoomMemberEdge):
			if !ok {
				return nil
			}
			return graphql.WriterFunc(func(w io.Writer) {
				w.Write([]byte{'{'})
				graphql.MarshalString(field.Alias).MarshalGQL(w)
				w.Write([]byte{':'})
				ec.marshalNRoomMemberEdge2ᚖjourneyhubᚋentᚐRoomMemberEdge(ctx, field.Selections, res).MarshalGQL(w)
				w.Write([]byte{'}'})
			})
		case <-ctx.Done():
			return nil
		}
	}
}

func (ec *executionContext) fieldContext_Subscription_roomMemberUpdated(_ context.Context, field graphql.CollectedField) (fc *graphql.FieldContext, err error) {
	fc = &graphql.FieldContext{
		Object:     "Subscription",
		Field:      field,
		IsMethod:   true,
		IsResolver: true,
		Child: func(ctx context.Context, field graphql.CollectedField) (*graphql.FieldContext, error) {
			switch field.Name {
			case "node":
				return ec.fieldContext_RoomMemberEdge_node(ctx, field)
			case "cursor":
				return ec.fieldContext_RoomMemberEdge_cursor(ctx, field)
			}
			return nil, fmt.Errorf("no field named %q was found under type RoomMemberEdge", field.Name)
		},
	}
	return fc, nil
}

func (ec *executionContext) _Subscription_roomMemberDeleted(ctx context.Context, field graphql.CollectedField) (ret func(ctx context.Context) graphql.Marshaler) {
	fc, err := ec.fieldContext_Subscription_roomMemberDeleted(ctx, field)
	if err != nil {
		return nil
	}
	ctx = graphql.WithFieldContext(ctx, fc)
	defer func() {
		if r := recover(); r != nil {
			ec.Error(ctx, ec.Recover(ctx, r))
			ret = nil
		}
	}()
	resTmp, err := ec.ResolverMiddleware(ctx, func(rctx context.Context) (interface{}, error) {
		ctx = rctx // use context from middleware stack in children
		return ec.resolvers.Subscription().RoomMemberDeleted(rctx)
	})
	if err != nil {
		ec.Error(ctx, err)
		return nil
	}
	if resTmp == nil {
		if !graphql.HasFieldError(ctx, fc) {
			ec.Errorf(ctx, "must not be null")
		}
		return nil
	}
	return func(ctx context.Context) graphql.Marshaler {
		select {
		case res, ok := <-resTmp.(<-chan pulid.ID):
			if !ok {
				return nil
			}
			return graphql.WriterFunc(func(w io.Writer) {
				w.Write([]byte{'{'})
				graphql.MarshalString(field.Alias).MarshalGQL(w)
				w.Write([]byte{':'})
				ec.marshalNID2journeyhubᚋentᚋschemaᚋpulidᚐID(ctx, field.Selections, res).MarshalGQL(w)
				w.Write([]byte{'}'})
			})
		case <-ctx.Done():
			return nil
		}
	}
}

func (ec *executionContext) fieldContext_Subscription_roomMemberDeleted(_ context.Context, field graphql.CollectedField) (fc *graphql.FieldContext, err error) {
	fc = &graphql.FieldContext{
		Object:     "Subscription",
		Field:      field,
		IsMethod:   true,
		IsResolver: true,
		Child: func(ctx context.Context, field graphql.CollectedField) (*graphql.FieldContext, error) {
			return nil, errors.New("field of type ID does not have child fields")
		},
	}
	return fc, nil
}

// endregion **************************** field.gotpl *****************************

// region    **************************** input.gotpl *****************************

func (ec *executionContext) unmarshalInputSendMessageInput(ctx context.Context, obj interface{}) (model.SendMessageInput, error) {
	var it model.SendMessageInput
	asMap := map[string]interface{}{}
	for k, v := range obj.(map[string]interface{}) {
		asMap[k] = v
	}

	fieldsInOrder := [...]string{"roomID", "notifyUserID", "replyTo", "content", "files", "voice", "links"}
	for _, k := range fieldsInOrder {
		v, ok := asMap[k]
		if !ok {
			continue
		}
		switch k {
		case "roomID":
			ctx := graphql.WithPathContext(ctx, graphql.NewPathWithField("roomID"))
			data, err := ec.unmarshalNID2journeyhubᚋentᚋschemaᚋpulidᚐID(ctx, v)
			if err != nil {
				return it, err
			}
			it.RoomID = data
		case "notifyUserID":
			ctx := graphql.WithPathContext(ctx, graphql.NewPathWithField("notifyUserID"))
			data, err := ec.unmarshalOID2ᚖjourneyhubᚋentᚋschemaᚋpulidᚐID(ctx, v)
			if err != nil {
				return it, err
			}
			it.NotifyUserID = data
		case "replyTo":
			ctx := graphql.WithPathContext(ctx, graphql.NewPathWithField("replyTo"))
			data, err := ec.unmarshalOID2ᚖjourneyhubᚋentᚋschemaᚋpulidᚐID(ctx, v)
			if err != nil {
				return it, err
			}
			it.ReplyTo = data
		case "content":
			ctx := graphql.WithPathContext(ctx, graphql.NewPathWithField("content"))
			data, err := ec.unmarshalOString2ᚖstring(ctx, v)
			if err != nil {
				return it, err
			}
			it.Content = data
		case "files":
			ctx := graphql.WithPathContext(ctx, graphql.NewPathWithField("files"))
			data, err := ec.unmarshalOUploadMessageFileInput2ᚕᚖjourneyhubᚋgraphᚋmodelᚐUploadMessageFileInputᚄ(ctx, v)
			if err != nil {
				return it, err
			}
			it.Files = data
		case "voice":
			ctx := graphql.WithPathContext(ctx, graphql.NewPathWithField("voice"))
			data, err := ec.unmarshalOUploadMessageVoiceInput2ᚖjourneyhubᚋgraphᚋmodelᚐUploadMessageVoiceInput(ctx, v)
			if err != nil {
				return it, err
			}
			it.Voice = data
		case "links":
			ctx := graphql.WithPathContext(ctx, graphql.NewPathWithField("links"))
			data, err := ec.unmarshalOCreateMessageLinkInput2ᚕᚖjourneyhubᚋgraphᚋmodelᚐCreateMessageLinkInputᚄ(ctx, v)
			if err != nil {
				return it, err
			}
			it.Links = data
		}
	}

	return it, nil
}

func (ec *executionContext) unmarshalInputUpdateMessageInput(ctx context.Context, obj interface{}) (model.UpdateMessageInput, error) {
	var it model.UpdateMessageInput
	asMap := map[string]interface{}{}
	for k, v := range obj.(map[string]interface{}) {
		asMap[k] = v
	}

	fieldsInOrder := [...]string{"content", "replaceLinks"}
	for _, k := range fieldsInOrder {
		v, ok := asMap[k]
		if !ok {
			continue
		}
		switch k {
		case "content":
			ctx := graphql.WithPathContext(ctx, graphql.NewPathWithField("content"))
			data, err := ec.unmarshalNString2string(ctx, v)
			if err != nil {
				return it, err
			}
			it.Content = data
		case "replaceLinks":
			ctx := graphql.WithPathContext(ctx, graphql.NewPathWithField("replaceLinks"))
			data, err := ec.unmarshalOCreateMessageLinkInput2ᚕᚖjourneyhubᚋgraphᚋmodelᚐCreateMessageLinkInputᚄ(ctx, v)
			if err != nil {
				return it, err
			}
			it.ReplaceLinks = data
		}
	}

	return it, nil
}

func (ec *executionContext) unmarshalInputUploadMessageFileInput(ctx context.Context, obj interface{}) (model.UploadMessageFileInput, error) {
	var it model.UploadMessageFileInput
	asMap := map[string]interface{}{}
	for k, v := range obj.(map[string]interface{}) {
		asMap[k] = v
	}

	fieldsInOrder := [...]string{"type", "file"}
	for _, k := range fieldsInOrder {
		v, ok := asMap[k]
		if !ok {
			continue
		}
		switch k {
		case "type":
			ctx := graphql.WithPathContext(ctx, graphql.NewPathWithField("type"))
			data, err := ec.unmarshalNMessageAttachmentType2journeyhubᚋentᚋmessageattachmentᚐType(ctx, v)
			if err != nil {
				return it, err
			}
			it.Type = data
		case "file":
			ctx := graphql.WithPathContext(ctx, graphql.NewPathWithField("file"))
			data, err := ec.unmarshalNUpload2githubᚗcomᚋ99designsᚋgqlgenᚋgraphqlᚐUpload(ctx, v)
			if err != nil {
				return it, err
			}
			it.File = data
		}
	}

	return it, nil
}

func (ec *executionContext) unmarshalInputUploadMessageVoiceInput(ctx context.Context, obj interface{}) (model.UploadMessageVoiceInput, error) {
	var it model.UploadMessageVoiceInput
	asMap := map[string]interface{}{}
	for k, v := range obj.(map[string]interface{}) {
		asMap[k] = v
	}

	fieldsInOrder := [...]string{"length", "file"}
	for _, k := range fieldsInOrder {
		v, ok := asMap[k]
		if !ok {
			continue
		}
		switch k {
		case "length":
			ctx := graphql.WithPathContext(ctx, graphql.NewPathWithField("length"))
			data, err := ec.unmarshalNUint642uint64(ctx, v)
			if err != nil {
				return it, err
			}
			it.Length = data
		case "file":
			ctx := graphql.WithPathContext(ctx, graphql.NewPathWithField("file"))
			data, err := ec.unmarshalNUpload2githubᚗcomᚋ99designsᚋgqlgenᚋgraphqlᚐUpload(ctx, v)
			if err != nil {
				return it, err
			}
			it.File = data
		}
	}

	return it, nil
}

// endregion **************************** input.gotpl *****************************

// region    ************************** interface.gotpl ***************************

// endregion ************************** interface.gotpl ***************************

// region    **************************** object.gotpl ****************************

var subscriptionImplementors = []string{"Subscription"}

func (ec *executionContext) _Subscription(ctx context.Context, sel ast.SelectionSet) func(ctx context.Context) graphql.Marshaler {
	fields := graphql.CollectFields(ec.OperationContext, sel, subscriptionImplementors)
	ctx = graphql.WithFieldContext(ctx, &graphql.FieldContext{
		Object: "Subscription",
	})
	if len(fields) != 1 {
		ec.Errorf(ctx, "must subscribe to exactly one stream")
		return nil
	}

	switch fields[0].Name {
	case "messageCreated":
		return ec._Subscription_messageCreated(ctx, fields[0])
	case "messageUpdated":
		return ec._Subscription_messageUpdated(ctx, fields[0])
	case "messageDeleted":
		return ec._Subscription_messageDeleted(ctx, fields[0])
	case "roomMemberCreated":
		return ec._Subscription_roomMemberCreated(ctx, fields[0])
	case "roomMemberUpdated":
		return ec._Subscription_roomMemberUpdated(ctx, fields[0])
	case "roomMemberDeleted":
		return ec._Subscription_roomMemberDeleted(ctx, fields[0])
	default:
		panic("unknown field " + strconv.Quote(fields[0].Name))
	}
}

// endregion **************************** object.gotpl ****************************

// region    ***************************** type.gotpl *****************************

func (ec *executionContext) unmarshalNSendMessageInput2journeyhubᚋgraphᚋmodelᚐSendMessageInput(ctx context.Context, v interface{}) (model.SendMessageInput, error) {
	res, err := ec.unmarshalInputSendMessageInput(ctx, v)
	return res, graphql.ErrorOnPath(ctx, err)
}

func (ec *executionContext) unmarshalNUpdateMessageInput2journeyhubᚋgraphᚋmodelᚐUpdateMessageInput(ctx context.Context, v interface{}) (model.UpdateMessageInput, error) {
	res, err := ec.unmarshalInputUpdateMessageInput(ctx, v)
	return res, graphql.ErrorOnPath(ctx, err)
}

func (ec *executionContext) unmarshalNUploadMessageFileInput2ᚖjourneyhubᚋgraphᚋmodelᚐUploadMessageFileInput(ctx context.Context, v interface{}) (*model.UploadMessageFileInput, error) {
	res, err := ec.unmarshalInputUploadMessageFileInput(ctx, v)
	return &res, graphql.ErrorOnPath(ctx, err)
}

func (ec *executionContext) unmarshalOUploadMessageFileInput2ᚕᚖjourneyhubᚋgraphᚋmodelᚐUploadMessageFileInputᚄ(ctx context.Context, v interface{}) ([]*model.UploadMessageFileInput, error) {
	if v == nil {
		return nil, nil
	}
	var vSlice []interface{}
	if v != nil {
		vSlice = graphql.CoerceList(v)
	}
	var err error
	res := make([]*model.UploadMessageFileInput, len(vSlice))
	for i := range vSlice {
		ctx := graphql.WithPathContext(ctx, graphql.NewPathWithIndex(i))
		res[i], err = ec.unmarshalNUploadMessageFileInput2ᚖjourneyhubᚋgraphᚋmodelᚐUploadMessageFileInput(ctx, vSlice[i])
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

func (ec *executionContext) unmarshalOUploadMessageVoiceInput2ᚖjourneyhubᚋgraphᚋmodelᚐUploadMessageVoiceInput(ctx context.Context, v interface{}) (*model.UploadMessageVoiceInput, error) {
	if v == nil {
		return nil, nil
	}
	res, err := ec.unmarshalInputUploadMessageVoiceInput(ctx, v)
	return &res, graphql.ErrorOnPath(ctx, err)
}

// endregion ***************************** type.gotpl *****************************
