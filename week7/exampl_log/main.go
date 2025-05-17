package main

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

func main() {
	//prodLog()
	//sampleLog()
	//sugarLog()
}

func prodLog() {
	logg := zap.Must(zap.NewProduction())

	defer logg.Sync()

	logg.Info("Hi its Zap!")

	logg.Info("P Log", zap.String("name", "Zap"), zap.Int("age", 5))
}

func sugarLog() {
	logg := zap.Must(zap.NewProduction())

	logg.Sugar().Infow("Sugar logger", "1234", "userID")
}

func sampleLog() {
	sto := zapcore.AddSync(os.Stdout)
	level := zap.NewAtomicLevelAt(zap.WarnLevel)
	prodCfg := zap.NewDevelopmentEncoderConfig()

	jEnc := zapcore.NewJSONEncoder(prodCfg)

	jOutCore := zapcore.NewCore(jEnc, sto, level)

	sampCore := zapcore.NewSamplerWithOptions(jOutCore, time.Second, 3, 0)

	log := zap.New(sampCore)

	for i := 0; i <= 10; i++ {
		log.Info("IT IS INFO")
		log.Warn("IT IS Warn")
		log.Error("IT IS Deb")
	}
}
