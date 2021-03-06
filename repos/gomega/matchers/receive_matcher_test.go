package matchers_test

import (
    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
    . "github.com/onsi/gomega/matchers"
    "time"
)

type kungFuActor interface {
    DrunkenMaster() bool
}

type jackie struct {
    name string
}

func (j *jackie) DrunkenMaster() bool {
    return true
}

var _ = Describe("ReceiveMatcher", func() {
    Context("with no argument", func() {
        Context("for a buffered channel", func() {
            It("should succeed", func() {
                channel := make(chan bool, 1)

                Ω(channel).ShouldNot(Receive())

                channel <- true

                Ω(channel).Should(Receive())
            })
        })

        Context("for an unbuffered channel", func() {
            It("should succeed (eventually)", func() {
                channel := make(chan bool)

                Ω(channel).ShouldNot(Receive())

                go func() {
                    time.Sleep(10 * time.Millisecond)
                    channel <- true
                }()

                Eventually(channel).Should(Receive())
            })
        })
    })

    Context("with a pointer argument", func() {
        Context("of the correct type", func() {
            It("should write the value received on the channel to the pointer", func() {
                channel := make(chan int, 1)

                var value int

                Ω(channel).ShouldNot(Receive(&value))
                Ω(value).Should(BeZero())

                channel <- 17

                Ω(channel).Should(Receive(&value))
                Ω(value).Should(Equal(17))
            })
        })

        Context("to various types of objects", func() {
            It("should work", func() {
                //channels of strings
                stringChan := make(chan string, 1)
                stringChan <- "foo"

                var s string
                Ω(stringChan).Should(Receive(&s))
                Ω(s).Should(Equal("foo"))

                //channels of slices
                sliceChan := make(chan []bool, 1)
                sliceChan <- []bool{true, true, false}

                var sl []bool
                Ω(sliceChan).Should(Receive(&sl))
                Ω(sl).Should(Equal([]bool{true, true, false}))

                //channels of channels
                chanChan := make(chan chan bool, 1)
                c := make(chan bool)
                chanChan <- c

                var receivedC chan bool
                Ω(chanChan).Should(Receive(&receivedC))
                Ω(receivedC).Should(Equal(c))

                //channels of interfaces
                jackieChan := make(chan kungFuActor, 1)
                aJackie := &jackie{name: "Jackie Chan"}
                jackieChan <- aJackie

                var theJackie kungFuActor
                Ω(jackieChan).Should(Receive(&theJackie))
                Ω(theJackie).Should(Equal(aJackie))
            })
        })

        Context("of the wrong type", func() {
            It("should error", func() {
                channel := make(chan int)
                var incorrectType bool

                success, _, err := (&ReceiveMatcher{Arg: &incorrectType}).Match(channel)
                Ω(success).Should(BeFalse())
                Ω(err).Should(HaveOccurred())

                var notAPointer int
                success, _, err = (&ReceiveMatcher{Arg: notAPointer}).Match(channel)
                Ω(success).Should(BeFalse())
                Ω(err).Should(HaveOccurred())
            })
        })
    })

    Context("When actual is a *closed* channel", func() {
        Context("for a buffered channel", func() {
            It("should work until it hits the end of the buffer", func() {
                channel := make(chan bool, 1)
                channel <- true

                close(channel)

                Ω(channel).Should(Receive())

                success, _, err := (&ReceiveMatcher{}).Match(channel)
                Ω(success).Should(BeFalse())
                Ω(err).Should(HaveOccurred())
            })
        })

        Context("for an unbuffered channel", func() {
            It("should error", func() {
                channel := make(chan bool)
                close(channel)

                success, _, err := (&ReceiveMatcher{}).Match(channel)
                Ω(success).Should(BeFalse())
                Ω(err).Should(HaveOccurred())
            })
        })
    })

    Context("When actual is a send-only channel", func() {
        It("should error", func() {
            channel := make(chan bool)

            var writerChannel chan<- bool
            writerChannel = channel

            success, _, err := (&ReceiveMatcher{}).Match(writerChannel)
            Ω(success).Should(BeFalse())
            Ω(err).Should(HaveOccurred())
        })
    })

    Context("when acutal is a non-channel", func() {
        It("should error", func() {
            var nilChannel chan bool

            success, _, err := (&ReceiveMatcher{}).Match(nilChannel)
            Ω(success).Should(BeFalse())
            Ω(err).Should(HaveOccurred())

            success, _, err = (&ReceiveMatcher{}).Match(nil)
            Ω(success).Should(BeFalse())
            Ω(err).Should(HaveOccurred())

            success, _, err = (&ReceiveMatcher{}).Match(3)
            Ω(success).Should(BeFalse())
            Ω(err).Should(HaveOccurred())
        })
    })
})
