package protocol

import (
	"github.com/lucas-clemente/quic-go"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Version", func() {
	isReservedVersion := func(v quic.VersionNumber) bool {
		return v&0x0f0f0f0f == 0x0a0a0a0a
	}

	It("says if a version is valid", func() {
		Expect(IsValidVersion(VersionTLS)).To(BeTrue())
		Expect(IsValidVersion(VersionWhatever)).To(BeFalse())
		Expect(IsValidVersion(VersionUnknown)).To(BeFalse())
		Expect(IsValidVersion(VersionDraft29)).To(BeTrue())
		Expect(IsValidVersion(Version1)).To(BeTrue())
		Expect(IsValidVersion(Version2)).To(BeTrue())
		Expect(IsValidVersion(1234)).To(BeFalse())
	})

	It("versions don't have reserved version numbers", func() {
		Expect(isReservedVersion(VersionTLS)).To(BeFalse())
	})

	It("has the right string representation", func() {
		Expect(VersionWhatever.String()).To(Equal("whatever"))
		Expect(VersionUnknown.String()).To(Equal("unknown"))
		Expect(VersionDraft29.String()).To(Equal("draft-29"))
		Expect(Version1.String()).To(Equal("v1"))
		Expect(Version2.String()).To(Equal("v2"))
		// check with unsupported version numbers from the wiki
		Expect(quic.VersionNumber(0x51303039).String()).To(Equal("gQUIC 9"))
		Expect(quic.VersionNumber(0x51303133).String()).To(Equal("gQUIC 13"))
		Expect(quic.VersionNumber(0x51303235).String()).To(Equal("gQUIC 25"))
		Expect(quic.VersionNumber(0x51303438).String()).To(Equal("gQUIC 48"))
		Expect(quic.VersionNumber(0x01234567).String()).To(Equal("0x1234567"))
	})

	It("recognizes supported versions", func() {
		Expect(IsSupportedVersion(SupportedVersions, 0)).To(BeFalse())
		Expect(IsSupportedVersion(SupportedVersions, SupportedVersions[0])).To(BeTrue())
		Expect(IsSupportedVersion(SupportedVersions, SupportedVersions[len(SupportedVersions)-1])).To(BeTrue())
	})

	Context("highest supported version", func() {
		It("finds the supported version", func() {
			supportedVersions := []quic.VersionNumber{1, 2, 3}
			other := []quic.VersionNumber{6, 5, 4, 3}
			ver, ok := ChooseSupportedVersion(supportedVersions, other)
			Expect(ok).To(BeTrue())
			Expect(ver).To(Equal(quic.VersionNumber(3)))
		})

		It("picks the preferred version", func() {
			supportedVersions := []quic.VersionNumber{2, 1, 3}
			other := []quic.VersionNumber{3, 6, 1, 8, 2, 10}
			ver, ok := ChooseSupportedVersion(supportedVersions, other)
			Expect(ok).To(BeTrue())
			Expect(ver).To(Equal(quic.VersionNumber(2)))
		})

		It("says when no matching version was found", func() {
			_, ok := ChooseSupportedVersion([]quic.VersionNumber{1}, []quic.VersionNumber{2})
			Expect(ok).To(BeFalse())
		})

		It("handles empty inputs", func() {
			_, ok := ChooseSupportedVersion([]quic.VersionNumber{102, 101}, []quic.VersionNumber{})
			Expect(ok).To(BeFalse())
			_, ok = ChooseSupportedVersion([]quic.VersionNumber{}, []quic.VersionNumber{1, 2})
			Expect(ok).To(BeFalse())
			_, ok = ChooseSupportedVersion([]quic.VersionNumber{}, []quic.VersionNumber{})
			Expect(ok).To(BeFalse())
		})
	})

	Context("reserved versions", func() {
		It("adds a greased version if passed an empty slice", func() {
			greased := GetGreasedVersions([]quic.VersionNumber{})
			Expect(greased).To(HaveLen(1))
			Expect(isReservedVersion(greased[0])).To(BeTrue())
		})

		It("creates greased lists of version numbers", func() {
			supported := []quic.VersionNumber{10, 18, 29}
			for _, v := range supported {
				Expect(isReservedVersion(v)).To(BeFalse())
			}
			var greasedVersionFirst, greasedVersionLast, greasedVersionMiddle int
			// check that
			// 1. the greased version sometimes appears first
			// 2. the greased version sometimes appears in the middle
			// 3. the greased version sometimes appears last
			// 4. the supported versions are kept in order
			for i := 0; i < 100; i++ {
				greased := GetGreasedVersions(supported)
				Expect(greased).To(HaveLen(4))
				var j int
				for i, v := range greased {
					if isReservedVersion(v) {
						if i == 0 {
							greasedVersionFirst++
						}
						if i == len(greased)-1 {
							greasedVersionLast++
						}
						greasedVersionMiddle++
						continue
					}
					Expect(supported[j]).To(Equal(v))
					j++
				}
			}
			Expect(greasedVersionFirst).ToNot(BeZero())
			Expect(greasedVersionLast).ToNot(BeZero())
			Expect(greasedVersionMiddle).ToNot(BeZero())
		})
	})
})