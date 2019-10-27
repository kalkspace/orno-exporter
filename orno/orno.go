package orno

import (
	"encoding/binary"
	"math"

	"github.com/goburrow/modbus"
	"github.com/kalkspace/orno-exporter/exporter"
	"github.com/kalkspace/orno-exporter/orno/models"
	"github.com/sirupsen/logrus"
)

type Meter interface {
	MetricRegisters() []models.RegisterConfig
}

var _ exporter.StateSource = (*Reader)(nil)

type Reader struct {
	client modbus.Client
	meter  Meter
	log    logrus.FieldLogger
}

func NewReader(logger logrus.FieldLogger, address string, meter Meter) *Reader {
	handler := modbus.NewRTUClientHandler(address)
	handler.BaudRate = 9600
	handler.DataBits = 8
	handler.Parity = "E"
	handler.StopBits = 1
	handler.SlaveId = 1
	handler.Logger = modbusLogger(logger)

	client := modbus.NewClient(handler)

	return &Reader{
		client: client,
		meter:  meter,
		log:    logger,
	}
}

func (r *Reader) Metrics() []string {
	configs := r.meter.MetricRegisters()
	metrics := make([]string, len(configs))
	for i, conf := range configs {
		metrics[i] = conf.Label
	}
	return metrics
}

func regLimits(meter Meter) (min uint16, max uint16) {
	min = math.MaxUint16
	max = 0

	for _, conf := range meter.MetricRegisters() {
		if conf.Num < min {
			min = conf.Num
		}
		end := conf.Num + uint16(conf.Size)
		if end > max {
			max = end
		}
	}

	return
}

const registerBulkSize = 64

func (r *Reader) Fetch() (map[string]interface{}, error) {
	min, max := regLimits(r.meter)
	log := r.log.WithField("min", min).WithField("max", max)
	log.Debug("Starting to read registers in bulk")

	registers := make([]byte, max*2) // registers are 16 bit (2 byte)
	for i := min; i < max; i += registerBulkSize {
		size := uint16(math.Min(float64(max-i), registerBulkSize))
		log := log.WithField("addr", i)
		log.WithField("size", size).Debug("Reading holding registers...")

		result, err := r.client.ReadHoldingRegisters(i, size)
		if err != nil {
			log.WithError(err).Warn("Reading failed")
			continue
			//return nil, err
		}
		copy(registers[i*2:], result)
	}
	log.WithField("registers", registers).Debug("bulk read finished")

	state := make(map[string]interface{})
	for _, conf := range r.meter.MetricRegisters() {
		val, ok := r.readFromConf(registers, conf)
		if ok {
			state[conf.Label] = val
		}
	}
	return state, nil
}

func (r *Reader) readFromConf(registers []byte, conf models.RegisterConfig) (interface{}, bool) {
	start := int(conf.Num * 2)
	end := start + conf.Size*2
	bytes := registers[start:end]

	log := r.log.WithFields(logrus.Fields{
		"bytes":  bytes,
		"metric": conf.Label,
	})

	var val interface{}
	switch conf.Type {
	case models.BinTypeUint16:
		val = binary.BigEndian.Uint16(bytes)
	case models.BinTypeFloat:
		bits := binary.BigEndian.Uint32(bytes)
		raw := math.Float32frombits(bits)
		val = math.Round(float64(raw)*100) / 100
	default:
		log.Debug("Discarding metric")
		return 0, false
	}

	log.WithField("value", val).Debug("Read metric")
	return val, true
}
