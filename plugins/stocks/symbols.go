package stocks

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-redis/redis"
	"gitlab.com/Cacophony/Processor/plugins/stocks/iex"
	"go.uber.org/zap"
)

const (
	symbolsExpiration = time.Hour * 24
)

var (
	regions = []string{"US", "DE"}
)

func symbolsKey(region string) string {
	return fmt.Sprintf("cacophony:processor:stocks:symbols:%s", region)
}

func (p *Plugin) getSymbolsForRegion(ctx context.Context, region string) ([]*iex.Symbol, error) {
	key := symbolsKey(region)
	var symbols []*iex.Symbol

	symbolsRaw, err := p.redis.Get(key).Bytes()
	if err == redis.Nil {
		symbols, err = p.iexClient.RefDataSymbolsInternational(ctx, region)
		if err != nil {
			return nil, err
		}

		symbolsRaw, err = json.Marshal(symbols)
		if err != nil {
			return nil, err
		}

		err = p.redis.Set(key, symbolsRaw, symbolsExpiration).Err()
		if err != nil {
			return nil, err
		}

		p.logger.Info("got Symbols from API, and cached",
			zap.String("region", region),
			zap.String("key", key),
			zap.Duration("expiration", symbolsExpiration),
		)

		return symbols, nil
	} else if err != nil {
		return nil, err
	}

	err = json.Unmarshal(symbolsRaw, &symbols)

	return symbols, err
}

func (p *Plugin) getAllSymbols(ctx context.Context) ([]*iex.Symbol, error) {
	var allSymbols []*iex.Symbol // nolint: prealloc
	for _, region := range regions {
		symbols, err := p.getSymbolsForRegion(ctx, region)
		if err != nil {
			return nil, err
		}

		allSymbols = append(allSymbols, symbols...)
	}

	return allSymbols, nil
}

func (p *Plugin) lookupSymbol(ctx context.Context, symbol string) (*iex.Symbol, error) {
	symbols, err := p.getAllSymbols(ctx)
	if err != nil {
		return nil, err
	}

	for _, symbolItem := range symbols {
		if strings.EqualFold(symbolItem.Symbol, symbol) {
			return symbolItem, nil
		}
	}

	return nil, errors.New("symbol not found")
}
