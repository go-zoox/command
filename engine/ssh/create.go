package ssh

import (
	"fmt"
	"net"

	sshx "golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

func (s *ssh) create() error {
	var err error

	var hostkeyCallback sshx.HostKeyCallback
	if s.cfg.KnowHostsFilePath != "" {
		hostkeyCallback, err = knownhosts.New(s.cfg.KnowHostsFilePath)
		if err != nil {
			return fmt.Errorf("failed to read known_hosts: %v", err)
		}
	}

	sshConf := &sshx.ClientConfig{
		User: s.cfg.User,
		// Auth:            []sshx.AuthMethod{},
		// HostKeyCallback: hostkeyCallback,
		HostKeyCallback: func(hostname string, remote net.Addr, key sshx.PublicKey) error {
			if s.cfg.IsIgnoreStrictHostKeyChecking {
				return nil
			}

			if hostkeyCallback == nil {
				return nil
			}

			return hostkeyCallback(hostname, remote, key)
		},
		HostKeyAlgorithms: []string{
			sshx.KeyAlgoED25519,
			sshx.KeyAlgoRSA,
			sshx.KeyAlgoDSA,
			sshx.KeyAlgoRSASHA512,
			sshx.KeyAlgoRSASHA256,
			sshx.CertAlgoRSASHA512v01,
			sshx.CertAlgoRSASHA256v01,
			sshx.CertAlgoRSAv01,
			sshx.CertAlgoDSAv01,
			sshx.CertAlgoECDSA256v01,
			sshx.CertAlgoECDSA384v01,
			sshx.CertAlgoECDSA521v01,
			sshx.CertAlgoED25519v01,
			sshx.KeyAlgoECDSA256,
			sshx.KeyAlgoECDSA384,
			sshx.KeyAlgoECDSA521,
		},
	}

	if s.cfg.PrivateKey != "" {
		var signer sshx.Signer
		if s.cfg.PrivateKeySecret == "" {
			signer, err = sshx.ParsePrivateKey([]byte(s.cfg.PrivateKey))
			if err != nil {
				if err.Error() != (&sshx.PassphraseMissingError{}).Error() {
					return fmt.Errorf("failed to parse private key: %v", err)
				}

				return err
			}
		} else {
			signer, err = sshx.ParsePrivateKeyWithPassphrase([]byte(s.cfg.PrivateKey), []byte(s.cfg.PrivateKeySecret))
			if err != nil {
				return fmt.Errorf("failed to parse private key (protected): %v", err)
			}
		}

		sshConf.Auth = append(sshConf.Auth, sshx.PublicKeys(signer))
	} else {
		if s.cfg.Pass != "" {
			sshConf.Auth = append(sshConf.Auth, sshx.Password(s.cfg.Pass))
		}
	}

	s.client, err = sshx.Dial("tcp", fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port), sshConf)
	if err != nil {
		return err
	}

	s.session, err = s.client.NewSession()
	if err != nil {
		return err
	}

	return nil
}
