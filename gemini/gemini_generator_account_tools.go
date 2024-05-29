package gemini

import (
	"encoding/json"
	"log/slog"
	"peekaboo_tools/safemap"
	"strings"
	"time"
)

func ParseGeminiAccountJson() {
	var gAccountInfo GlobalAccountInfo
	err := json.Unmarshal([]byte(defaultConfigStr), &gAccountInfo)
	if err != nil {
		panic(err)
	}
	var newAccountInfo []*AccountInfo
	split := strings.Split(newAccountInfoStr, "\n")

	for _, item := range split {
		if strings.TrimSpace(item) == "" {
			continue
		}

		item = item[len(item)-35:]
		newAccountInfo = append(newAccountInfo, &AccountInfo{
			SecurityInfo: &SecurityInfo{Token: item},
		})
	}
	gAccountInfo.AccountInfos = newAccountInfo

	indent, err := json.MarshalIndent(gAccountInfo, "", "  ")
	if err != nil {
		panic(err)
	}
	slog.Info(string(indent))
}

type (
	GlobalAccountInfo struct {
		AccountInfos             []*AccountInfo                     `json:"accountInfos"`
		AccountExpiredInfo       *AccountExpiredInfo                `json:"commonAccountExpiredInfo"`
		RequestRateLimitingInfos []*RequestRateLimitingInfo         `json:"commonRequestRateLimitingInfos"`
		RequestTokenLimitInfos   []*RequestTokenLimitInfo           `json:"commonRequestTokenLimitInfos"`
		AccountInfoMap           *safemap.Map[string, *AccountInfo] `json:"-"`
	}

	//AccountInfo 账户基本信息
	AccountInfo struct {
		SecurityInfo             *SecurityInfo              `json:"securityInfo"`
		AccountExpiredInfo       *AccountExpiredInfo        `json:"accountExpiredInfo,omitempty"`
		RequestRateLimitingInfos []*RequestRateLimitingInfo `json:"requestRateLimitingInfos,omitempty"`
		RequestTokenLimitInfos   []*RequestTokenLimitInfo   `json:"requestTokenLimitInfos,omitempty"`
	}

	//SecurityInfo 安全信息
	SecurityInfo struct {
		Token   string    `json:"token"`
		ErrType ErrorType `json:"-"`
	}

	//RequestRateLimitingInfo 请求速率限制
	RequestRateLimitingInfo struct {
		RequestRateLimiting              int           `json:"request_rate_limiting"` //请求限速
		RequestRateLimitingLimitTimeUnit LimitTimeUnit `json:"request_rate_unit"`     //限制时间单位 1 秒 2分钟  3 时 4天 5 月 6季度 7 年
	}
	//RequestTokenLimitInfo token 令牌限制
	RequestTokenLimitInfo struct {
		RequestTokenLimiting              int           `json:"request_token_limiting"` //令牌限制数量
		RequestTokenLimitingLimitTimeUnit LimitTimeUnit `json:"request_token_unit"`     //限制时间单位 1 秒 2分钟  3 时 4天 5 月 6季度 7 年
	}

	//AccountExpiredInfo 账户过期信息
	AccountExpiredInfo struct {
		AccountExpiredAtStr string    `json:"account_expired_at"` //过期限制
		AccountExpiredAt    time.Time `json:"-"`                  //过期限制
	}

	NotAvailableAccountInfo struct {
		Token string `json:"token"`
	}
)

// LimitTimeUnit 限制时间单位
type LimitTimeUnit int

const (
	LimitTimeUnitSecond LimitTimeUnit = 1
	LimitTimeUnitMinute LimitTimeUnit = 2
	LimitTimeUnitHour   LimitTimeUnit = 3
	LimitTimeUnitDay    LimitTimeUnit = 4
	LimitTimeUnitMonth  LimitTimeUnit = 5
	LimitTimeUnitSeason LimitTimeUnit = 6
	LimitTimeUnitYear   LimitTimeUnit = 7
)

type ErrorType int

const (
	ErrorTypeNotAvailable     ErrorType = 1  //不可用
	ErrorTypeRequestRateLimit ErrorType = 2  //限速
	ErrorTypeTokenLimit       ErrorType = 3  //令牌限制
	ErrorTypeExpired          ErrorType = 4  //过期
	ErrorTypeAccountError     ErrorType = 5  //账号错误
	ErrorTypeNotFoundError    ErrorType = 98 //账号错误
	ErrorTypeNotOtherInfo     ErrorType = 99 //其他限制
)

var newAccountInfoStr = `
neillgathings310@gmail.com----5ohnl0tzoi----ezreilya@hotmail.com-----AIzaSyB9o9St-VqzzF9iCXtckKxfNL-2AL8-AYI
wingrectrippduchsjessica@gmail.com----cNe5Nk6C3G3F----renrelisbio@email.com---AIzaSyAr5E2fsguF_f6ZtjbzGdXX9CEe5JHrTzM
signmatcgreenmortaaron@gmail.com----vxAY9KLNeF0R----vermyonulpai@outlook.cz----AIzaSyCH635qHm0Nk9Sxbz4Gm4WoCEEdwCRtPJE
vildicaltuashley@gmail.com----JmjI6mSBMnm----truthewipov@outlook.com.br-----AIzaSyCehkVprLm3bLueJv6IRphQjo-n9gdnPNQ
updacmanejim@gmail.com----90shUn9D7JRc----mestcongtanggen@outlook.co.il----AIzaSyD4FE0MyV04foCxy0BUvS1rmsjK1mZRxyM
miljanovicin@gmail.com----DMedJWom50IY----exicenwah@mailbox.org----AIzaSyB2j8OjvtyJA2J7MLD2xBZgchHEsRJPA4E
icitdebu1977@gmail.com----vFpivHmENwPQTJ----spotimunas@mailbox.org----AIzaSyDFI4FtinV2sKlXye3lMCJh-lnQ2c-07eA
finneganmike050@gmail.com----36L2cfnH5Ggb----ddevubrysli@mail.com----AIzaSyD7AVtUQx5hExMmY7sNhYYr9YtH5q_M3u4
sinpesthemo1972@gmail.com----fsjV1bay8TolR----swapadlogleu@mailbox.org----AIzaSyACnnFUiwZlt7bfV9PuE6V7QOxEHdbxZ20
wohlmansam1@gmail.com----Fyk5lcpeBGrnW----dentiemali@mail.com----AIzaSyDYD63Ql3keE0f_GMVK5zUHECbbKZuwALM
`

var defaultConfigStr = `
{
  "accountInfos": [
    {
      "securityInfo": {
        "token": "AIzaSyB9o9St-VqzzF9iCXtckKxfNL-2AL8-AYI"
      }
    },
    {
      "securityInfo": {
        "token": "AIzaSyAr5E2fsguF_f6ZtjbzGdXX9CEe5JHrTzM"
      }
    },
    {
      "securityInfo": {
        "token": "AIzaSyCH635qHm0Nk9Sxbz4Gm4WoCEEdwCRtPJE"
      }
    },
    {
      "securityInfo": {
        "token": "AIzaSyCehkVprLm3bLueJv6IRphQjo-n9gdnPNQ"
      }
    },
    {
      "securityInfo": {
        "token": "AIzaSyD4FE0MyV04foCxy0BUvS1rmsjK1mZRxyM"
      }
    },
    {
      "securityInfo": {
        "token": "AIzaSyB2j8OjvtyJA2J7MLD2xBZgchHEsRJPA4E"
      }
    },
    {
      "securityInfo": {
        "token": "AIzaSyDFI4FtinV2sKlXye3lMCJh-lnQ2c-07eA"
      }
    },
    {
      "securityInfo": {
        "token": "AIzaSyD7AVtUQx5hExMmY7sNhYYr9YtH5q_M3u4"
      }
    },
    {
      "securityInfo": {
        "token": "AIzaSyACnnFUiwZlt7bfV9PuE6V7QOxEHdbxZ20"
      }
    },
    {
      "securityInfo": {
        "token": "AIzaSyDYD63Ql3keE0f_GMVK5zUHECbbKZuwALM"
      }
    }
  ],
  "commonAccountExpiredInfo": {
  },
  "commonRequestRateLimitingInfos": [
    {
      "request_rate_limiting": 1500,
      "request_rate_unit": 4
    },
    {
      "request_rate_limiting": 15,
      "request_rate_unit": 2
    }
  ],
  "commonRequestTokenLimitInfos": [
    {
      "request_token_limiting": 32000,
      "request_token_unit": 2
    }
  ]
}
`
