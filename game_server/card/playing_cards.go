package card

type PlayingCards struct {
	cardsInHand			[]*Cards		//手上的牌
	cardsAlreadyChi		[]*Cards		//已经吃了的牌
	cardsAlreadyPeng	[]*Cards		//已经碰了的牌
	cardsAlreadyGang	[]*Cards		//已经杠了的牌

	cardsForCheckHu		*Cards			//用于临时检查是否胡牌
}

func NewPlayingCards() *PlayingCards {
	cards :=  &PlayingCards{
		cardsForCheckHu:	NewCards(),
	}
	cards.cardsInHand = cards.initCardsSlice()
	cards.cardsAlreadyChi = cards.initCardsSlice()
	cards.cardsAlreadyPeng = cards.initCardsSlice()
	cards.cardsAlreadyGang = cards.initCardsSlice()
	return cards
}

//增加一张牌
func (playingCards *PlayingCards) AddCard(card *Card) {
	playingCards.cardsInHand[card.CardType].AddAndSort(card)
}

//丢弃一张牌
func (playingCards *PlayingCards) DropCard(card *Card) bool {
	return playingCards.cardsInHand[card.CardType].TakeWay(card)
}

//吃牌，要吃whatCard，以及吃哪个组合whatGroup
func (playingCards *PlayingCards) Chi(whatCard *Card, whatGroup *Cards) bool {
	if !playingCards.canChi(whatCard, whatGroup) {
		return false
	}

	for _, card := range whatGroup.Data() {//移动除了whatCard以外的card到cardsAlreadyChi
		if card.SameAs(whatCard) {
			continue
		}
		playingCards.cardsInHand[card.CardType].TakeWay(card)
		playingCards.cardsAlreadyChi[whatCard.CardType].AppendCard(card)
	}

	//最后把whatCard加入cardsAlreadyChi
	playingCards.cardsAlreadyChi[whatCard.CardType].AddAndSort(whatCard)

	return true
}

//碰牌
func (playingCards *PlayingCards) Peng(whatCard *Card) bool {
	if !playingCards.canPeng(whatCard) {
		return false
	}

	playingCards.cardsInHand[whatCard.CardType].TakeWay(whatCard)
	playingCards.cardsInHand[whatCard.CardType].TakeWay(whatCard)
	playingCards.cardsAlreadyPeng[whatCard.CardType].AppendCard(whatCard)
	playingCards.cardsAlreadyPeng[whatCard.CardType].AppendCard(whatCard)
	playingCards.cardsAlreadyPeng[whatCard.CardType].AddAndSort(whatCard)
	return true
}

//杠牌
func (playingCards *PlayingCards) Gang(whatCard *Card) bool {
	if !playingCards.canGang(whatCard) {
		return false
	}

	playingCards.cardsInHand[whatCard.CardType].TakeWay(whatCard)
	playingCards.cardsInHand[whatCard.CardType].TakeWay(whatCard)
	playingCards.cardsInHand[whatCard.CardType].TakeWay(whatCard)
	playingCards.cardsAlreadyGang[whatCard.CardType].AppendCard(whatCard)
	playingCards.cardsAlreadyGang[whatCard.CardType].AppendCard(whatCard)
	playingCards.cardsAlreadyGang[whatCard.CardType].AppendCard(whatCard)
	playingCards.cardsAlreadyGang[whatCard.CardType].AddAndSort(whatCard)
	return true
}

func (playingCards *PlayingCards) ToString() string{
	str := ""
	str += "cardsInHand:\n" + playingCards.cardsSliceToString(playingCards.cardsInHand)
	str += "cardsAlreadyChi:\n" + playingCards.cardsSliceToString(playingCards.cardsAlreadyChi)
	str += "cardsAlreadyPeng:\n" + playingCards.cardsSliceToString(playingCards.cardsAlreadyPeng)
	str += "cardsAlreadyGang:\n" + playingCards.cardsSliceToString(playingCards.cardsAlreadyGang)
	return str
}

/*	计算指定的牌可以吃牌的组合
*/
func (playingCards *PlayingCards) ComputeChiGroup(card *Card) []*Cards {
	return playingCards.cardsInHand[card.CardType].computeChiGroup(card)
}

//检查是否能吃
func (playingCards *PlayingCards) canChi(whatCard *Card, whatGroup *Cards) bool {
	return playingCards.cardsInHand[whatCard.CardType].canChi(whatCard, whatGroup)
}

//检查是否能碰
func (playingCards *PlayingCards) canPeng(whatCard *Card) bool  {
	return playingCards.cardsInHand[whatCard.CardType].canPeng(whatCard)
}

//检查是否能杠
func (playingCards *PlayingCards) canGang(whatCard *Card) bool {
	return playingCards.cardsInHand[whatCard.CardType].canGang(whatCard)
}


//初始化cards
func (playingCards *PlayingCards) initCardsSlice()[]*Cards {
	cardsSlice := make([]*Cards, Max_CardType)
	for idx := 0; idx < Max_CardType; idx++ {
		cardsSlice[idx] = NewCards()
	}
	return cardsSlice
}

func (playingCards *PlayingCards) cardsSliceToString(cardsSlice []*Cards) string{
	str := ""
	for _, cards := range cardsSlice{
		str += cards.ToString() + "\n"
	}
	return str
}

func (playingCards *PlayingCards) GetInHandCards(cardType int) *Cards{
	if cardType < 0 || cardType >= Max_CardType {
		return nil
	}
	return playingCards.cardsInHand[cardType]
}

func (playingCards *PlayingCards) GetAlreadyChiCards(cardType int) *Cards{
	if cardType < 0 || cardType >= Max_CardType {
		return nil
	}
	return playingCards.cardsAlreadyChi[cardType]
}

func (playingCards *PlayingCards) GetAlreadyPengCards(cardType int) *Cards{
	if cardType < 0 || cardType >= Max_CardType {
		return nil
	}
	return playingCards.cardsAlreadyPeng[cardType]
}

func (playingCards *PlayingCards) GetAlreadyGangCards(cardType int) *Cards{
	if cardType < 0 || cardType >= Max_CardType {
		return nil
	}
	return playingCards.cardsAlreadyGang[cardType]
}

func (playingCards *PlayingCards) IsHu() bool {
	playingCards.cardsForCheckHu.Clear()
	for _, cards := range playingCards.cardsInHand {
		mod := cards.Len() % 3
		if mod != 0 && mod != 2 {//每一种类型的牌的数量mod 3 不是 0或者2的话肯定不可能胡
			return false
		}
		playingCards.cardsForCheckHu.AppendCards(cards)
	}

	playingCards.cardsForCheckHu.Sort()
	return playingCards.cardsForCheckHu.IsHu()
}