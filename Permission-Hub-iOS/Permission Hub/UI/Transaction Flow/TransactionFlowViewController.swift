//
//  TransactionFlowViewController.swift
//  Permission Hub
//
//  Created by Corné on 13/04/2019.
//  Copyright © 2019 Planet. All rights reserved.
//

import UIKit

final class TransactionFlowViewController: PHTableViewController {

    // MARK: - Private properties

    private lazy var bottomStackView: UIStackView = {

        let stackView = UIStackView()
        stackView.translatesAutoresizingMaskIntoConstraints = false

        stackView.axis = .horizontal
        stackView.alignment = .center
        stackView.distribution = .equalCentering
        stackView.spacing = 20

        return stackView
    }()

    // MARK: - Life cycle

    init(transaction: TransactionNotification) {

        var items = [PHTableViewViewCellType]()
        items.append(.notification(
            type: .notification,
            text: "Verified"))
        items.append(.notification(
            type: .warning,
            text: "Permission warning!"))
        items.append(.description(
            title: transaction.title,
            description: transaction.description))

        transaction.items.forEach {
            items.append(.transactionItem(item: $0))
        }

        super.init(items: items)
    }

    required init?(coder aDecoder: NSCoder) {
        fatalError("init(coder:) has not been implemented")
    }

    override func viewDidLoad() {
        super.viewDidLoad()

        view.addSubview(bottomStackView)
        bottomStackView.heightAnchor.constraint(equalToConstant: 75).isActive = true
        bottomStackView.topAnchor.constraint(equalTo: super.tableView.bottomAnchor).isActive = true
        bottomStackView.centerXAnchor.constraint(equalTo: view.centerXAnchor).isActive = true
        bottomStackView.bottomAnchor.constraint(equalTo: view.safeAreaLayoutGuide.bottomAnchor).isActive = true

        for i in 0..<2 {

            let button = UIButton()
            button.translatesAutoresizingMaskIntoConstraints = false

            button.widthAnchor.constraint(equalToConstant: 92).isActive = true
            button.heightAnchor.constraint(equalToConstant: 34).isActive = true

            button.titleLabel?.font = PHFonts.regular()

            let color = i == 0 ? PHColors.greyishBrown : PHColors.topaz
            button.setTitleColor(color, for: .normal)
            button.backgroundColor = .white

            button.layer.cornerRadius = 5
            button.layer.borderWidth = 1
            button.layer.borderColor = color.cgColor

            let title = i == 0 ? "Decline" : "Continue"
            button.setTitle(title, for: .normal)

            if i == 0 {
                button.addTarget(
                    self,
                    action: #selector(declineButtonTapped),
                    for: .touchUpInside)
            } else {
                button.addTarget(
                    self,
                    action: #selector(continueButtonTapped),
                    for: .touchUpInside)
            }

            bottomStackView.addArrangedSubview(button)
        }
    }

    // MARK: - Networking

    private func respondToTransaction(isAccepted: Bool) {

        do {
            let service = NetworkingService.shared
            try service.respondToTransaction(withId: "123", isAccepted: true) { [unowned self] response in

                switch response {
                case .success:
                    self.navigationController?.popViewController(animated: true)
                case .failure(let error):
                    print(error)
                }
            }
        } catch {
            print(error)
        }
    }

    // MARK: - Selectors

    @objc private func declineButtonTapped(_ sender: UIButton) {
        respondToTransaction(isAccepted: false)
    }

    @objc private func continueButtonTapped(_ sender: UIButton) {
        respondToTransaction(isAccepted: true)
    }
}
