package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// WalletsFromJsonToMAP read json filename and path
func WalletsFromJsonToMAP(ctx context.Context, path, filename string) ([]map[string]interface{}, error) {
	_, span := otel.Tracer(nameSourcerJson).Start(ctx, "WalletsFromJsonToMAP")
	defer span.End()
	jsonFile, err := os.Open(fmt.Sprintf("%s/%s", path, filename))
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)

	var result []map[string]interface{}
	err = json.Unmarshal([]byte(byteValue), &result)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	span.SetAttributes(attribute.Int("jsonfile.query.result", len(result)))

	return result, nil
}

// WalletsFromJsonToStruct read json filename and path
func WalletsFromJsonToStruct(ctx context.Context, path, filename string) ([]Wallet, error) {
	_, span := otel.Tracer(nameSourcerJson).Start(ctx, "WalletsFromJsonToStruct")
	defer span.End()
	jsonFile, err := os.Open(fmt.Sprintf("%s/%s", path, filename))
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)

	var result []Wallet
	err = json.Unmarshal([]byte(byteValue), &result)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	span.SetAttributes(attribute.Int("jsonfile.query.result", len(result)))

	return result, nil
}
