package handle

import (
	"jiacrontab/libs/mailer"
	"jiacrontab/libs/proto"
	"jiacrontab/model"
	"jiacrontab/server/conf"
	"log"
)

type Logic struct{}

func (l *Logic) Register(args model.Client, reply *proto.MailArgs) error {

	*reply = proto.MailArgs{
		Host: conf.MailService.Host,
		User: conf.MailService.User,
		Pass: conf.MailService.Passwd,
	}

	ret := model.DB().Model(&model.Client{}).Where("addr=?", args.Addr).Update(&args)

	if ret.RowsAffected == 0 {
		ret = model.DB().Create(&args)
	}

	log.Println("register client", args)
	return ret.Error
}

func (l *Logic) Depends(args model.DependsTasks, reply *bool) error {
	log.Printf("Callee Logic.Depend taskId %s", args[0].TaskId)
	*reply = true
	for _, v := range args {
		if err := rpcCall(v.Dest, "Task.ExecDepend", v, &reply); err != nil {
			*reply = false
			return err
		}
	}

	return nil
}

func (l *Logic) DependDone(args model.DependsTask, reply *bool) error {
	log.Printf("Callee Logic.DependDone task %s", args.Name)
	*reply = true
	if err := rpcCall(args.Dest, "Task.ResolvedDepends", args, &reply); err != nil {
		*reply = false
		return err
	}

	return nil
}

func (l *Logic) SendMail(args proto.SendMail, reply *bool) error {
	mailer.SendMail(args.MailTo, args.Subject, args.Content)
	*reply = true
	return nil
}
