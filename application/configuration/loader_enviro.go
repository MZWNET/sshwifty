// Sshwifty - A Web SSH client
//
// Copyright (C) 2019 Rui NI <nirui@gmx.com>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package configuration

import (
	"fmt"
	"os"
	"strconv"

	"github.com/niruix/sshwifty/application/log"
	"github.com/niruix/sshwifty/application/network"
)

const (
	enviroTypeName = "Environment Variable"
)

// Enviro creates an environment variable based configuration loader
func Enviro() Loader {
	return func(log log.Logger) (string, Configuration, error) {
		log.Info("Loading configuration from environment variables ...")

		cfg := fileCfgCommon{
			HostName:       os.Getenv("SSHWIFTY_HOSTNAME"),
			SharedKey:      os.Getenv("SSHWIFTY_SHAREDKEY"),
			Socks5:         os.Getenv("SSHWIFTY_SOCKS5"),
			Socks5User:     os.Getenv("SSHWIFTY_SOCKS5_USER"),
			Socks5Password: os.Getenv("SSHWIFTY_SOCKS5_PASSWORD"),
		}

		listenPort, listenPortErr := strconv.ParseUint(
			os.Getenv("SSHWIFTY_LISTENPORT"), 10, 16)

		if listenPortErr != nil {
			return enviroTypeName, Configuration{}, fmt.Errorf(
				"Invalid \"SSHWIFTY_LISTENPORT\": %s", listenPortErr)
		}

		initialTimeout, _ := strconv.ParseUint(
			os.Getenv("SSHWIFTY_INITIALTIMEOUT"), 10, 32)

		readTimeout, _ := strconv.ParseUint(
			os.Getenv("SSHWIFTY_READTIMEOUT"), 10, 32)

		writeTimeout, _ := strconv.ParseUint(
			os.Getenv("SSHWIFTY_WRITETIMEOUT"), 10, 32)

		heartbeatTimeout, _ := strconv.ParseUint(
			os.Getenv("SSHWIFTY_HEARTBEATTIMEOUT"), 10, 32)

		readDelay, _ := strconv.ParseUint(
			os.Getenv("SSHWIFTY_READDELAY"), 10, 32)

		writeDelay, _ := strconv.ParseUint(
			os.Getenv("SSHWIFTY_WRITEELAY"), 10, 32)

		cfgSer := fileCfgServer{
			ListenInterface:       os.Getenv("SSHWIFTY_LISTENINTERFACE"),
			ListenPort:            uint16(listenPort),
			InitialTimeout:        int(initialTimeout),
			ReadTimeout:           int(readTimeout),
			WriteTimeout:          int(writeTimeout),
			HeartbeatTimeout:      int(heartbeatTimeout),
			ReadDelay:             int(readDelay),
			WriteDelay:            int(writeDelay),
			TLSCertificateFile:    os.Getenv("SSHWIFTY_TLSCERTIFICATEFILE"),
			TLSCertificateKeyFile: os.Getenv("SSHWIFTY_TLSCERTIFICATEKEYFILE"),
		}

		var dialer network.Dial

		if len(cfg.Socks5) <= 0 {
			dialer = network.TCPDial()
		} else {
			sDial, sDialErr := network.BuildSocks5Dial(
				cfg.Socks5, cfg.Socks5User, cfg.Socks5Password)

			if sDialErr != nil {
				return enviroTypeName, Configuration{}, sDialErr
			}

			dialer = sDial
		}

		return enviroTypeName, Configuration{
			HostName:  cfg.HostName,
			SharedKey: cfg.SharedKey,
			Dialer:    dialer,
			Servers:   []Server{cfgSer.build()},
		}, nil
	}
}